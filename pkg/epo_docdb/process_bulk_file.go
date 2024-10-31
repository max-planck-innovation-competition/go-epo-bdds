package epo_docdb

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/krolaw/zipstream"
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

const logLevelProcess = 5

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
			logger.With("country", c).Debug("skipping file")
			return true
		} else {
			logger.With("country", c).Debug("including file")
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
		return true // skip if file type not matched
	}
	return false // include if no file types are specified
}

// ContentHandler is a function that handles the content of a file
type ContentHandler func(fileName string, fileContent string)

// regexFileName is used to extract the filename by using attributes from the xml file
var regexFileName = regexp.MustCompile(`country="([A-Z]{1,3})".*doc-number="([A-Z0-9]{1,15})".*kind="([A-Z0-9]{1,3})"`)

// ProcessDirectory processes a directory
func (p *Processor) ProcessDirectory(workingDirectoryPath string) (err error) {
	directoryLogger := slog.With("wd", workingDirectoryPath)
	directoryLogger.Info("process directory")

	filePaths := []string{}
	// read the bulk zip file
	err = fs.WalkDir(os.DirFS(workingDirectoryPath), ".", func(path string, d fs.DirEntry, err error) error {
		// check if dir
		if d.IsDir() {
			return nil
		}
		// check if zip file and starts with "docdb_"
		if strings.Contains(path, ".zip") && strings.HasPrefix(path, "docdb_") {
			filePath := filepath.Join(workingDirectoryPath, path)
			filePaths = append(filePaths, filePath)
		}
		// default (other files)
		return nil
	})
	if err != nil {
		directoryLogger.With("err", err).Error("failed to walk dir")
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
			directoryLogger.With("filePath", filePath).Info("skipping file based on file type")
			continue
		}

		// add to queueFiles
		queueFiles = append(queueFiles, filePath)
	}

	for i, filePath := range queueFiles {
		directoryLogger.With("file", filePath).Info("processing file")
		// process bulk zip file
		err = p.ProcessBulkZipFile(filePath)
		if err != nil {
			directoryLogger.With("err", err).Error("failed to process bulk zip file")
			return err
		}
		// log the current progress
		directoryLogger.
			With("file", i+1).
			With("total", len(queueFiles)).
			Info("current progress")
	}

	directoryLogger.Info("successfully done")
	return

}

// ProcessBulkZipFile processes a bulk zip file
func (p *Processor) ProcessBulkZipFile(filePath string) (err error) {
	logger := slog.With("filePath", filePath)

	// Open the bulk zip file
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		logger.With("err", err).Error("failed to open bulk zip file")
		return err
	}
	defer reader.Close()

	queueFiles := []*zip.File{}

	// Iterate over the files in the zip archive
	for _, f := range reader.File {
		path := f.Name

		// check if dir
		if f.FileInfo().IsDir() {
			continue
		}
		// check if zip file
		if strings.Contains(path, "Root/DOC/") && strings.Contains(path, ".zip") {

			// skip countries that are not in the list of countries to include
			if len(p.includeAuthorities) > 0 {
				if p.skipFileBasedOnAuthority(path) {
					continue
				}
			}

			// check if state handler is set
			// if yes then check if the file is already done
			if p.StateHandler != nil {
				bulkState, _ := p.StateHandler.RegisterOrSkipZipFile(path)
				if bulkState == state_handler.Done {
					// if already done, skip
					logger.With("zipFile", path).Debug("skipping zip file")
					continue
				}
			}

			// add to queueFiles
			queueFiles = append(queueFiles, f)
		}
	}

	// Set the number of workers
	fileCh := make(chan *zip.File, len(queueFiles)) // Buffered channel with the number of files
	var wg sync.WaitGroup
	total := len(queueFiles)

	// Start the worker pool
	for w := 0; w < p.Workers; w++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for zipFile := range fileCh {

				workerLogger := slog.With("workerId", workerId).With("file", zipFile.Name)

				// process zip file
				p.ProcessZipFile(workerLogger, zipFile)

				// mark zip file as finished
				if p.StateHandler != nil {
					p.StateHandler.MarkZipFileAsFinished()
				}

				// log the current progress
				workerLogger.
					With("todo", len(fileCh)).
					With("total", total).
					Debug("worker processed zip file")
			}
		}(w)
	}

	// Send files to the workers
	for _, zipFile := range queueFiles {
		fileCh <- zipFile
	}
	close(fileCh)

	// Wait for all workers to finish
	wg.Wait()

	logger.Debug("successfully done")
	return
}

// ProcessZipFile processes a zip file within a bulk zip file
func (p *Processor) ProcessZipFile(logger *slog.Logger, zipFile *zip.File) {
	logger = logger.With("zipFile", zipFile.Name)

	// Open the zip file
	f, err := zipFile.Open()
	if err != nil {
		logger.With("err", err).Error("failed to open zip file")
		return
	}
	defer f.Close()

	// Use zipstream to process the zip entries without loading the entire file into memory
	zr := zipstream.NewReader(f)

	for {
		header, err := zr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.With("err", err).Error("failed to read zip entry")
			return
		}
		logger.With("xmlFile", header.Name).Debug("child found")

		// check state handler
		if p.StateHandler != nil {
			// check if the file is already done
			xmlStatus, _ := p.StateHandler.RegisterOrSkipXMLFile(header.Name, "/Root/DOC/")
			if xmlStatus == state_handler.Done {
				// if already done, skip
				logger.Debug("skipping xml file")
				continue
			}
		}

		// process zip file content
		err = p.ProcessZipFileContent(logger, header, zr)
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
func (p *Processor) ProcessZipFileContent(logger *slog.Logger, header *zip.FileHeader, zr *zipstream.Reader) (err error) {
	logger = logger.With("xmlFile", header.Name)
	logger.Debug("process xml file")

	// zr is already positioned at the file content
	// zr implements io.Reader for the file content

	return p.ProcessExchangeFileContent(logger, zr)
}

// ExchangeDocument represents the structure of the exchange-document
type ExchangeDocument struct {
	XMLName   xml.Name `xml:"exchange-document"`
	Country   string   `xml:"country,attr"`
	DocNumber string   `xml:"doc-number,attr"`
	Kind      string   `xml:"kind,attr"`
	InnerXML  string   `xml:",innerxml"`
}

// FileName constructs the file name from the document attributes
func (doc *ExchangeDocument) FileName() string {
	return fmt.Sprintf("%s-%s-%s.xml", doc.Country, doc.DocNumber, doc.Kind)
}

func extractFileName(line string) string {
	// Regular expression to extract country, doc-number, and kind attributes
	matches := regexFileName.FindStringSubmatch(line)
	if len(matches) == 4 {
		country := matches[1]
		docNumber := matches[2]
		kind := matches[3]
		return fmt.Sprintf("%s-%s-%s.xml", country, docNumber, kind)
	}
	// If attributes are not found, handle the error as needed
	return "unknown.xml"
}

// regex for line break
var regexLineBreak = regexp.MustCompile(`[\r\n]+`)

// ProcessExchangeFileContent processes an exchange file content
func (p *Processor) ProcessExchangeFileContent(logger *slog.Logger, fc io.Reader) (err error) {
	// scan file
	scanner := bufio.NewScanner(fc)
	// set the max capacity of the scanner
	const maxCapacity = 10 * 1024 * 1024 // 10 MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	// custom line break
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
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
		if strings.Contains(line, "<exch:exchange-document") {
			lineContent.Reset()
			lineContent.WriteString(line)
			// Extract fileName from the line
			fileName = extractFileName(line)
		} else if strings.Contains(line, "</exch:exchange-document>") {
			lineContent.WriteString(line)
			p.ContentHandler(fileName, lineContent.String())
			lineContent.Reset()
			fileName = ""
		} else {
			lineContent.WriteString(line)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.With("err", err).Error("scanner error")
		return err
	}

	logger.Debug("done with file")
	return nil
}
