package epo_docdb

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Processor creates a
type Processor struct {
	ContentHandler   ContentHandler
	includeCountries map[string]struct{}
}

// NewProcessor creates a new processor
// the default handler is PrintLineHandler
func NewProcessor() *Processor {
	p := Processor{
		ContentHandler: PrintLineHandler,
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

// IncludeAuthorities sets the authorities to include
// if no countries are included all authorities are included.
// This is useful if you only want to include e.g. data from the EPO
func (p *Processor) IncludeAuthorities(cs ...string) {
	p.includeCountries = map[string]struct{}{}
	for _, c := range cs {
		c = strings.ToUpper(c)
		p.includeCountries[c] = struct{}{}
	}
}

// ContentHandler is a function that handles the content of a file
type ContentHandler func(fileName, fileContent string)

// regexFileName is used to extract the filename from the xml file
var regexFileName = regexp.MustCompile(`country="([A-Z]{1,3})".*doc-number="([A-Z0-9]{1,15})".*kind="([A-Z0-9]{1,3})".*doc-id="([A-Z0-9]{1,20})"`)

// ProcessDirectory processes a directory
func (p *Processor) ProcessDirectory(workingDirectoryPath string) (err error) {
	logger := log.WithField("workingDirectoryPath", workingDirectoryPath)
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
		logger.WithError(err).Error("failed to walk dir")
		return err
	}
	// order files ascending
	sort.Strings(filePaths)

	// iterate over files
	for _, filePath := range filePaths {
		err = p.ProcessBulkZipFile(filePath)
		if err != nil {
			logger.WithError(err).Error("failed to process bulk zip file")
			return err
		}
	}

	logger.Info("successfully done")
	return

}

// ProcessBulkZipFile processes a bulk zip file
func (p *Processor) ProcessBulkZipFile(filePath string) (err error) {
	logger := log.WithField("filePath", filePath)
	logger.Info("start reading file")

	// read the bulk zip file
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		logger.WithError(err).Error("failed to open bulk zip file")
		return err
	}
	err = fs.WalkDir(reader, ".", func(path string, d fs.DirEntry, err error) error {
		// check if dir
		if d.IsDir() {
			return nil
		}
		// check if zip file
		if strings.Contains(path, "Root/DOC/") && strings.Contains(path, ".zip") {

			// skip countries that are not in the list of countries to include
			if len(p.includeCountries) > 0 {
				// get file Name e.g. DOCDB-202402-CreateDelete-PubDate20240105AndBefore-AR-0001.zip
				var countryRegex = regexp.MustCompile("-([A-Z]{2})-[0-9]{1,10}\\.zip")
				fileName := filepath.Base(path)
				// check if the file name contains a country
				country := countryRegex.FindStringSubmatch(fileName)
				if len(country) == 2 {
					c := strings.ToUpper(country[1])
					// check if the country is in the list of countries to include
					if _, ok := p.includeCountries[c]; !ok {
						// skip this file
						logger.WithField("country", c).Info("skipping file")
						return nil
					} else {
						logger.WithField("country", c).Info("including file")
					}
				}
			}

			f, errOpen := reader.Open(path)
			if errOpen != nil {
				err = errOpen
				logger.WithError(err).Error("failed to open file")
				return err
			}
			logger.WithField("zipFile", path).Info("found zip file")
			p.processZipFile(logger, f)
		}
		// default (other files)
		return nil
	})
	if err != nil {
		logger.WithError(err).Error("failed to walk dir")
		return err
	}
	// close
	err = reader.Close()
	if err != nil {
		logger.WithError(err).Error("failed to close bulk zip file")
		return err
	}

	logger.Info("successfully done")
	return
}

// processZipFile processes a bulk zip file
func (p *Processor) processZipFile(logger *log.Entry, f fs.File) {
	stats, _ := f.Stat()
	logger = logger.WithField("zipFile", stats.Name())
	// read file
	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Fatal(err)
	}

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		logger.WithField("xmlFile", zipFile.Name).Info("child found")
		err = p.processZipFileContent(logger, zipFile)
		if err != nil {
			logger.WithError(err).Error("failed to process zip file content")
			return
		}
	}

}

// processZipFileContent processes a zip file content
func (p *Processor) processZipFileContent(logger *log.Entry, file *zip.File) (err error) {
	logger = log.WithField("xmlFile", file.Name)
	logger.Info("process xml file")
	ctx := context.TODO()
	fc, err := file.Open()
	if err != nil {
		msg := "failed to open zip %s for reading: %s"
		err = fmt.Errorf(msg, file.Name, err)
		logger.Error(err)
		return
	}
	defer func() {
		errClose := fc.Close()
		if errClose != nil {
			ctx.Done()
			logger.Fatalf("Failed to close file: %s", errClose)
		}
	}()
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
					logger.Error(err)
					return
				}
				if len(regexExtractionResults) != 1 && len(regexExtractionResults[0]) != 1 {
					msg := "failed extract filename"
					err = fmt.Errorf(msg)
					logger.Error(err)
					return
				}
				fileName = regexExtractionResults[0][1] + "-" + regexExtractionResults[0][2] + "-" + regexExtractionResults[0][3] + "_" + regexExtractionResults[0][4] + ".xml"

				lineContent.WriteString(line)
				// if the line also contains the end of the file
				if strings.Contains(line, "</exch:exchange-document>") {
					p.ContentHandler(fileName, lineContent.String())
					lineContent.Reset()
					fileName = ""
					continue
				}
			} else {
				log.WithField("line", line).Error("failed to split line")
			}
		} else {
			// if the line contains the end of the file
			if strings.Contains(line, "</exch:exchange-document>") {
				lineContent.WriteString(line)
				p.ContentHandler(fileName, lineContent.String())
				lineContent.Reset()
				fileName = ""
				continue
			}
			lineContent.WriteString(line)
		}
	}
	logger.Info("done with file")
	return
}
