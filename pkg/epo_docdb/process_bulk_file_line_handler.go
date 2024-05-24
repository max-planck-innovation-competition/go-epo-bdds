package epo_docdb

import (
	"fmt"
	"github.com/max-planck-innovation-competition/go-epo-bdds/pkg/state_handler"
	"log/slog"
	"os"
	"path/filepath"
)

// FileExporterLineHandler processes the files and saves them in the destination folder
func FileExporterLineHandler(destinationFolderPath string) ContentHandler {
	logger := slog.With("destinationFolderPath", destinationFolderPath)
	logger.Debug("started")

	// check if destination folder exists
	if _, err := os.Stat(destinationFolderPath); os.IsNotExist(err) {
		// create folder
		err = os.MkdirAll(destinationFolderPath, os.ModePerm)
		if err != nil {
			logger.With("err", err).Error("failed to create destination folder")
			panic(err)
		}
		logger.Info("created destination folder")
	}

	return func(
		fileName string,
		fileContent string,
		recorder state_handler.StateHandler, //dummy
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
		slog.With("err", err).With("fileName", fileName).Error("can not open file")
		return
	}
	// write the data to the file
	_, errWrite := file.WriteString(fileContent)
	if errWrite != nil {
		slog.With("err", errWrite).With("fileName", fileName).Error("failed to write to buffer")
		return
	}
	// close the file
	errClose := file.Close()
	if errClose != nil {
		slog.With("err", errClose).With("fileName", fileName).Error("failed to close file")
		return
	}
}

// PrintLineHandler is a dummy handler that prints the file name and the file content
func PrintLineHandler(
	fileName string,
	fileContent string,
	recorder state_handler.StateHandler,
) {
	fmt.Println(fileName, fileContent)
}
