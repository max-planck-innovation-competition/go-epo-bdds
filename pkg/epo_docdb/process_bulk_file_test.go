package epo_docdb

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessBulkZipFile(t *testing.T) {
	// skipTest(t)
	log.SetLevel(log.TraceLevel)
	ass := assert.New(t)
	err := ProcessBulkZipFile(
		"./test-data/docdb_xml_202344_CreateDelete_001.zip", "./test-data/xml")
	ass.NoError(err)
}

func TestProcessBulkZipBackfile(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	ass := assert.New(t)
	dir := "./test-data/backfile"
	files, _ := os.ReadDir(dir)
	path, _ := filepath.Abs(dir)
	for _, file := range files {
		bulkfilepath := filepath.Join(path, file.Name())
		err := ProcessBulkZipFile(
			bulkfilepath, "./test-data/backfile/xml")
		ass.NoError(err)
	}
}

func TestProcessBulkZipfrontfile(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	ass := assert.New(t)

	dir := "./test-data/frontfile"
	files, _ := os.ReadDir(dir)
	path, _ := filepath.Abs(dir)
	for _, file := range files {
		bulkfilepath := filepath.Join(path, file.Name())
		err := ProcessBulkZipFile(
			bulkfilepath, "./test-data/frontfile/xml")
		ass.NoError(err)
	}
}
