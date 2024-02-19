package epo_docdb

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// FileExporterLineHandler processes the files and saves them in the destination folder
func FileExporterLineHandler(destinationFolderPath string) ContentHandler {
	logger := log.WithField("handler", "FileExporterLineHandler")
	logger.Trace("started")

	// check if destination folder exists
	if _, err := os.Stat(destinationFolderPath); os.IsNotExist(err) {
		// create folder
		err = os.MkdirAll(destinationFolderPath, os.ModePerm)
		if err != nil {
			logger.WithError(err).Error("failed to create destination folder")
			panic(err)
		}
		logger.Info("created destination folder")
	}

	return func(
		fileName string,
		fileContent string,
	) {
		// join path
		filePath := filepath.Join(destinationFolderPath, fileName)
		SaveFile(filePath, fileContent)
	}
}

// SaveFile saves the file with the given content and the given name
func SaveFile(
	fileName string,
	fileContent string,
) {
	// create the file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("can not create fileName", err)
		return
	}
	// write the data to the file
	_, errWrite := file.WriteString(fileContent)
	if errWrite != nil {
		msg := "failed to write to buffer: %s"
		log.Fatalf(msg, errWrite)
		return
	}
	// close the file
	errClose := file.Close()
	if errClose != nil {
		msg := "failed to write close file: %s"
		log.Fatalf(msg, errClose)
		return
	}
}

// PrintLineHandler is a dummy handler that prints the file name and the file content
func PrintLineHandler(
	fileName string,
	fileContent string,
) {
	fmt.Println(fileName, fileContent)
}
