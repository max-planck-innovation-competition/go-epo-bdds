package docdb

import (
	"archive/zip"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

func ProcessBulkFile(sourceFile, destinationFolder string) (err error) {
	logger := log.WithField("zipFile", sourceFile)
	logger.Info("reading file")
	// read zip file
	readCloser, err := zip.OpenReader(sourceFile)
	if err != nil {
		msg := "Failed to open: %s"
		logger.Fatalf(msg, err)
	}
	logger.Info("file read")
	// close file after read
	defer func() {
		errClose := readCloser.Close()
		if errClose != nil {
			logger.Fatalf("Failed to close file: %s", errClose)
		}
	}()
	logger.Info("iterate over files")
	// iterate over all files in the zip directory
	for _, file := range readCloser.File {
		logger.WithField("filename", file.Name).Info("found")
		extension := filepath.Ext(file.Name)
		extension = strings.ToLower(extension)
		if extension == ".xml" {
			if errFile := processZippedFiles(file, destinationFolder); errFile != nil {
				return
			}
		} else {
			logger.WithField("filename", file.Name).Info("skipping file")
		}
	}
	logger.Info("successfully done")
	return
}

func processZippedFiles(file *zip.File, destinationFolder string) (err error) {
	logger := log.WithField("filename", file.Name).WithField("routine", "main")
	logger.Info("process file")
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
	return
}
