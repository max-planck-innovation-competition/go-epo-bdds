package epo_docdb

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/max-planck-innovation-competition/go-epo-bdds/pkg/state_handler"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Processor creates a
type Processor struct {
	ContentHandler     ContentHandler              // content handler
	includeAuthorities map[string]struct{}         // e.g. EP, WO, etc.
	includeFileTypes   map[string]struct{}         // e.g. CreateDelete, Amend, etc.
	StateHandler       *state_handler.StateHandler // optional state handler
	Workers            int                         // number of workers
}

// NewProcessor creates a new processor
// the default handler is PrintLineHandler
func NewProcessor() *Processor {
	p := Processor{
		ContentHandler: PrintLineHandler,
		Workers:        1,
	}
	return &p
}

// NewFileExportProcessor creates a new processor
// the default handler is FileExporterLineHandler
func NewFileExportProcessor(destinationPath string) *Processor {
	handler := FileExporterLineHandler(destinationPath)
	p := Processor{
		ContentHandler: handler,
	}
	return &p
}

// SetContentHandler sets the content handler
// you can create your own ContentHandler
func (p *Processor) SetContentHandler(fn ContentHandler) *Processor {
	p.ContentHandler = fn
	return p
}

// SetStateHandler adds a state handler
func (p *Processor) SetStateHandler(stateHandler *state_handler.StateHandler) *Processor {
	p.StateHandler = stateHandler
	return p
}

// IncludeAuthorities sets the authorities to include
// if no countries are included all authorities are included.
// This is useful if you only want to include e.g. data from the EPO
func (p *Processor) IncludeAuthorities(cs ...string) {
	p.includeAuthorities = map[string]struct{}{}
	for _, c := range cs {
		c = strings.ToUpper(c)
		p.includeAuthorities[c] = struct{}{}
	}
}

// skipFileBasedOnAuthority checks if the file should be skipped
// based on the authority
func (p *Processor) skipFileBasedOnAuthority(filePath string) bool {
	logger := slog.With("filePath", filePath)
	// get file Name e.g. DOCDB-202402-CreateDelete-PubDate20240105AndBefore-AR-0001.zip
	var countryRegex = regexp.MustCompile("-([A-Z]{2})-[0-9]{1,10}\\.zip")
	fileName := filepath.Base(filePath)
	// check if the file name contains a country
	country := countryRegex.FindStringSubmatch(fileName)
	if len(country) == 2 {
		c := strings.ToUpper(country[1])
		// check if the country is in the list of countries to include
		if _, ok := p.includeAuthorities[c]; !ok {
			// skip this file
			logger.With("country", c).Info("skipping file")
			return true
		} else {
			logger.With("country", c).Info("including file")
			return false
		}
	}
	logger.Warn("could not extract country from file name")
	return true // skip
}

// IncludeFileTypes sets the file types to include
// if no file types are included all file types are included.
// This is useful if you only want to include e.g. CreateDelete or Amend files
func (p *Processor) IncludeFileTypes(cs ...string) {
	p.includeFileTypes = map[string]struct{}{}
	for _, c := range cs {
		c = strings.ToUpper(c)
		p.includeFileTypes[c] = struct{}{}
	}
}

// skipFileBasedOnFileType checks if the file should be skipped
// based on the file type.
// e.g. CreateDelete, Amend, etc.
func (p *Processor) skipFileBasedOnFileType(filePath string) bool {
	// check if file types are included
	if len(p.includeFileTypes) > 0 {
		// iterate over file types
		for fileType := range p.includeFileTypes {
			// check if the file type is in the path
			if strings.Contains(strings.ToLower(filePath), strings.ToLower(fileType)) {
				return false
			}
		}
	}
	return true
}

// ContentHandler is a function that handles the content of a file
type ContentHandler func(fileName, fileContent string)

// regexFileName is used to extract the filename by using attributes from the xml file
var regexFileName = regexp.MustCompile(`country="([A-Z]{1,3})".*doc-number="([A-Z0-9]{1,15})".*kind="([A-Z0-9]{1,3})"`)

// ProcessDirectory processes a directory
func (p *Processor) ProcessDirectory(workingDirectoryPath string) (err error) {
	logger := slog.With("workingDirectoryPath", workingDirectoryPath)
	logger.Info("start reading file")

	filePaths := []string{}
	// read the bulk zip file
	err = fs.WalkDir(os.DirFS(workingDirectoryPath), ".", func(path string, d fs.DirEntry, err error) error {
		// check if dir
		if d.IsDir() {
			return nil
		}
		// check if zip file and starts with "doc_db"
		if strings.Contains(path, ".zip") && strings.HasPrefix(path, "docdb_") {
			filePath := filepath.Join(workingDirectoryPath, path)
			filePaths = append(filePaths, filePath)
		}
		// default (other files)
		return nil
	})
	if err != nil {
		logger.With("err", err).Error("failed to walk dir")
		return err
	}
	// order files ascending
	sort.Strings(filePaths)

	queueFiles := []string{}
	// iterate over files
	for _, filePath := range filePaths {
		// check if state handler is set
		if p.StateHandler != nil {
			// check if the file is already done
			state, _ := p.StateHandler.RegisterOrSkipZipFile(filePath)
			if state == state_handler.Done {
				// if already done, skip
				continue
			}
		}
		// skip file based on file type
		if p.skipFileBasedOnFileType(filePath) {
			logger.With("filePath", filePath).Info("skipping file based on file type")
			continue
		}

		// add to queueFiles
		queueFiles = append(queueFiles, filePath)
	}

	for i, filePath := range queueFiles {
		// process bulk zip file
		err = p.ProcessBulkZipFile(filePath)
		if err != nil {
			logger.With("err", err).Error("failed to process bulk zip file")
			return err
		}
		// log the current progress
		logger.
			With("file", i+1).
			With("total", len(queueFiles)).
			Info("processed file")
	}

	logger.Info("successfully done")
	return

}

// ProcessBulkZipFile processes a bulk zip file
func (p *Processor) ProcessBulkZipFile(filePath string) (err error) {
	logger := slog.With("filePath", filePath)
	logger.Info("start reading file")

	// read the bulk zip file
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		logger.With("err", err).Error("failed to open bulk zip file")
		return err
	}

	queueFiles := []string{}

	err = fs.WalkDir(reader, ".", func(path string, d fs.DirEntry, err error) error {
		// check if dir
		if d.IsDir() {
			return nil
		}
		// check if zip file
		if strings.Contains(path, "Root/DOC/") && strings.Contains(path, ".zip") {

			// skip countries that are not in the list of countries to include
			if len(p.includeAuthorities) > 0 {
				if p.skipFileBasedOnAuthority(path) {
					return nil
				}
			}

			// check if state handler is set
			// if yes then check if the file is already done
			if p.StateHandler != nil {
				bulkState, _ := p.StateHandler.RegisterOrSkipZipFile(path)
				if bulkState == state_handler.Done {
					// if already done, skip
					logger.With("zipFile", path).Info("skipping zip file")
					return nil
				}
			}

			// add to queueFiles
			queueFiles = append(queueFiles, path)
		}
		// default (other files)
		return nil
	})
	if err != nil {
		logger.With("err", err).Error("failed to walk dir")
		return err
	}

	// Set the number of workers
	numWorkers := 5
	fileCh := make(chan string, len(queueFiles)) // Buffered channel with the number of files
	var wg sync.WaitGroup
	total := len(queueFiles)

	// Start the worker pool
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for path := range fileCh {

				workerLogger := slog.With("workerId", workerId).With("file", path)

				// open file
				f, errOpen := reader.Open(path)
				if errOpen != nil {
					workerLogger.With("err", errOpen).Error("failed to open file")
					continue
				}

				workerLogger.Info("worker processing zip file")
				// process zip file
				p.ProcessZipFile(logger, f)

				// mark zip file as finished
				if p.StateHandler != nil {
					p.StateHandler.MarkZipFileAsFinished()
				}

				errClose := f.Close()
				if errClose != nil {
					workerLogger.With("err", errClose).Error("failed to close file")
					return
				} // Ensure the file is closed after processing

				// log the current progress
				workerLogger.
					With("todo", len(fileCh)).
					With("total", total).
					Info("worker processed zip file")
			}
		}(w)
	}

	// Send files to the workers
	for _, path := range queueFiles {
		fileCh <- path
	}
	close(fileCh) // Close the channel to signal workers that no more files will be sent

	// Wait for all workers to finish
	wg.Wait()

	// close
	err = reader.Close()
	if err != nil {
		logger.With("err", err).Error("failed to close bulk zip file")
		return err
	}

	logger.Info("successfully done")
	return
}

// ProcessZipFile processes a bulk zip file
func (p *Processor) ProcessZipFile(logger *slog.Logger, f fs.File) {
	stats, _ := f.Stat()
	logger = logger.With("zipFile", stats.Name())
	// read file
	data, err := io.ReadAll(f)
	if err != nil {
		logger.With("err", err).Error("failed to read zip file")
		return
	}
	// create a new zip reader
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		logger.With("err", err).Error("failed to read zip file")
		return
	}

	// read all the files from zip archive
	for _, zipFile := range zipReader.File {
		logger.With("xmlFile", zipFile.Name).Info("child found")
		// check state handler
		if p.StateHandler != nil {
			// check if the file is already done
			xmlStatus, _ := p.StateHandler.RegisterOrSkipXMLFile(zipFile.Name, "/Root/DOC/")
			if xmlStatus == state_handler.Done {
				// if already done, skip
				logger.Info("skipping xml file")
				continue
			}
		}
		// process zip file content
		err = p.ProcessZipFileContent(logger, zipFile)
		if err != nil {
			logger.With("err", err).Error("failed to process zip file content")
			return
		}
		// mark xml as finished
		if p.StateHandler != nil {
			p.StateHandler.MarkXMLAsFinished()
		}
	}

}

// ProcessZipFileContent processes a zip file content
func (p *Processor) ProcessZipFileContent(logger *slog.Logger, file *zip.File) (err error) {
	logger = logger.With("xmlFile", file.Name)
	logger.Info("process xml file")
	ctx := context.TODO()
	fc, err := file.Open()
	if err != nil {
		msg := "failed to open zip %s for reading: %s"
		err = fmt.Errorf(msg, file.Name, err)
		logger.With("err", err).Error("failed to open zip file")
		return
	}
	defer func() {
		errClose := fc.Close()
		if errClose != nil {
			ctx.Done()
			logger.With("err", errClose).Error("Failed to close file")
		}
	}()
	return p.ProcessExchangeFileContent(logger, fc)
}

// ProcessExchangeFileContent processes a exchange file content
func (p *Processor) ProcessExchangeFileContent(logger *slog.Logger, fc io.Reader) (err error) {
	// scan file
	scanner := bufio.NewScanner(fc)
	// set the max capacity of the scanner
	const maxCapacity = 500 * 1024 * 1024 // 500 MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	// custom line break
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		// regex for line break
		var regexLineBreak = regexp.MustCompile(`[\r\n]+`)
		loc := regexLineBreak.FindIndex(data)
		if len(loc) == 0 {
			return 0, nil, nil
		}
		i := loc[0]
		if i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})

	var lineContent strings.Builder
	var fileName string

	for scanner.Scan() {
		line := scanner.Text()
		// last line
		const docStart = "<exch:exchange-document "
		// start of file e.g. first line
		if strings.Contains(line, docStart) {
			// split the line
			split := strings.Split(line, docStart)
			if len(split) == 2 {
				line = docStart + split[1]
				// extract the filename
				regexExtractionResults := regexFileName.FindAllStringSubmatch(line, -1)
				if len(regexExtractionResults) == 0 {
					msg := "failed extract filename"
					err = fmt.Errorf(msg)
					logger.With("line", line).With("err", err).Error("failed to extract filename")
					return
				}
				if len(regexExtractionResults) != 1 && len(regexExtractionResults[0]) != 1 {
					msg := "failed extract filename"
					err = fmt.Errorf(msg)
					logger.With("line", line).With("err", err).Error("failed to extract filename")
					return
				}
				fileName = regexExtractionResults[0][1] + "-" + regexExtractionResults[0][2] + "-" + regexExtractionResults[0][3] + ".xml"

				lineContent.WriteString(line)
				// if the line also contains the end of the file
				if strings.Contains(line, "</exch:exchange-document>") {
					p.ContentHandler(fileName, lineContent.String())
					lineContent.Reset()
					fileName = ""
					if p.StateHandler != nil {
						p.StateHandler.MarkExchangeFileAsFinished()
					}
					continue
				}
			} else {
				slog.With("line", line).Error("failed to split line")
			}
		} else {
			// if the line contains the end of the file
			if strings.Contains(line, "</exch:exchange-document>") {
				lineContent.WriteString(line)
				p.ContentHandler(fileName, lineContent.String())
				lineContent.Reset()
				fileName = ""
				if p.StateHandler != nil {
					p.StateHandler.MarkExchangeFileAsFinished()
				}
				continue
			}
			lineContent.WriteString(line)
		}
	}
	logger.Info("done with file")
	return
}
