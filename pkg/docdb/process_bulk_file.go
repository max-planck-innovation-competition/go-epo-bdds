package docdb

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
	"regexp"
	"strings"
	"sync"
)

// regexFileName is used to extract the filename from the xml file
var regexFileName = regexp.MustCompile(`country="([A-Z]{1,3})".*doc-number="([A-Z0-9]{1,15})".*kind="([A-Z0-9]{1,3})".*doc-id="([A-Z0-9]{1,20})"`)

// ProcessBulkZipFile processes a bulk zip file
func ProcessBulkZipFile(bulkZipFile, destinationFolder string) (err error) {
	logger := log.WithField("bulkZipFile", bulkZipFile)
	logger.Info("start reading file")

	// read the bulk zip file
	reader, err := zip.OpenReader(bulkZipFile)
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
			f, errOpen := reader.Open(path)
			if errOpen != nil {
				err = errOpen
				logger.WithError(err).Error("failed to open file")
				return err
			}
			logger.WithField("zipFile", path).Info("found zip file")
			processZipFile(logger, f, destinationFolder)
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
func processZipFile(logger *log.Entry, f fs.File, destinationFolder string) {
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
		err = processZipFileContent(logger, zipFile, destinationFolder)
		if err != nil {
			logger.WithError(err).Error("failed to process zip file content")
			return
		}
	}

}

// processZipFileContent processes a zip file content
func processZipFileContent(logger *log.Entry, file *zip.File, destinationFolder string) (err error) {
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
	// init channels and sync
	// init channels and sync
	var wg sync.WaitGroup
	chContent := make(chan string)
	chFilename := make(chan string)
	chFileEnd := make(chan bool)
	// start 2nd process
	go fileWriter(ctx, destinationFolder, &wg, chContent, chFilename, chFileEnd)
	// scan file
	// scan file
	scanner := bufio.NewScanner(fc)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024*200) // 200 MB
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
				if len(regexExtractionResults) != 1 && len(regexExtractionResults[0]) != 1 {
					msg := "failed extract filename"
					err = fmt.Errorf(msg)
					logger.Error(err)
					return
				}
				filename := regexExtractionResults[0][1] + "-" + regexExtractionResults[0][2] + "-" + regexExtractionResults[0][3] + "_" + regexExtractionResults[0][4] + ".xml"
				wg.Add(1)
				chFilename <- filename
				// add the content
				wg.Add(1)
				chContent <- line
				// if the line also contains the end of the file
				if strings.Contains(line, "</exch:exchange-document>") {
					wg.Add(1)
					chFileEnd <- true
				}
			} else {
				log.WithField("line", line).Error("failed to split line")
			}
		} else {
			if strings.Contains(line, "</exch:exchange-document>") {
				// end of the file
				wg.Add(1)
				chContent <- line
				wg.Add(1)
				chFileEnd <- true
				continue
			}
			// normal line
			wg.Add(1)
			chContent <- line
		}
	}

	logger.Info("done with file")

	return
}

// fileWriter writes the content to a file
func fileWriter(
	ctx context.Context,
	destinationFolder string,
	wg *sync.WaitGroup,
	chContent <-chan string,
	chFilename <-chan string,
	chFileEnd <-chan bool,
) {
	logger := log.WithField("routine", "writer")
	logger.Trace("started")
	filename := ""
	var buf strings.Builder

	for {
		select {
		case <-ctx.Done():
			logger.Info("received context done")
			return
		case end := <-chFileEnd:
			if end {
				// if the last string was transmitted
				logger.Trace("last stuff was transmitted")
				// create a new file based on the filename
				if len(filename) == 0 {
					msg := "failed to extract filename: %s"
					logger.Fatalf(msg, filename)
					return
				}
				file, err := os.Create(destinationFolder + "/" + filename)
				if err != nil {
					logger.Error(err)
					return
				}
				// write the data to the file
				_, errWrite := file.WriteString(buf.String())
				if errWrite != nil {
					ctx.Done()
					msg := "failed to write to buffer: %s"
					logger.Fatalf(msg, errWrite)
					return
				}
				// close the file
				errClose := file.Close()
				if errClose != nil {
					ctx.Done()
					msg := "failed to write close file: %s"
					logger.Fatalf(msg, errClose)
					return
				}
				// clear the string builder
				buf.Reset()
				// clear the filename
				filename = ""
				break
			}
		case content := <-chContent:
			logger.
				// WithField("content", content).
				Trace("received data")
			// skip the empty line if there is nothing in the buffer
			if buf.Len() == 0 && len(content) == 0 {
				break
			}
			// if there is content write it to the buffer
			_, errWrite := buf.WriteString(content + "\n")
			if errWrite != nil {
				ctx.Done()
				msg := "failed to write to buffer: %s"
				logger.Fatalf(msg, errWrite)
				return
			}
			break
		case filename = <-chFilename:
			logger.WithField("filename", filename).Debug("Set filename")
			break
		}
		log.Trace("done")
		wg.Done()
	}

}
