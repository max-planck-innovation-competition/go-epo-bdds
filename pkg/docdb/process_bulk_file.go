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
	"strings"
)

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
			f, _ := reader.Open(path)
			logger.WithField("zipFile", path).Info("found zip file")
			processZipFile(logger, f)
		}
		// default (other files)
		return nil
	})
	// close
	err = reader.Close()
	if err != nil {
		logger.WithError(err).Error("failed to close bulk zip file")
		return err
	}

	logger.Info("successfully done")
	return
}

func processZipFile(logger *log.Entry, f fs.File) {
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
		processZipFileContent(logger, zipFile)
	}

}

func processZipFileContent(logger *log.Entry, file *zip.File) (err error) {
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
	// scan file
	reader := bufio.NewReader(fc)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		fmt.Println(string(line))
		// last line
		if strings.Contains(string(line), "</exch:exchange-document>") {

		}
		// start of file e.g. first line
		if strings.Contains(string(line), "<exch:exchange-document") {

		}
	}

	logger.Info("done with file")

	return
}
