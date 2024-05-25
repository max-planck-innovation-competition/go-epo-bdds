package epo_docdb

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProcessBulkZipFile(t *testing.T) {
	// skipTest(t)
	p := NewProcessor()
	err := p.ProcessBulkZipFile("./test-data/docdb_xml_202402_CreateDelete_001.zip")
	if err != nil {
		t.Error(err)
	}
}

func TestProcessBulkZipFile2023(t *testing.T) {
	// skipTest(t)
	ass := assert.New(t)
	p := NewFileExportProcessor("./test-data/xml")
	err := p.ProcessBulkZipFile("./test-data/docdb_xml_202402_CreateDelete_001.zip")
	ass.NoError(err)
}

func TestProcessEpFiles2023(t *testing.T) {
	// skipTest(t)
	ass := assert.New(t)
	p := NewFileExportProcessor("./test-data/xml/eps")
	p.IncludeAuthorities("EP")
	err := p.ProcessBulkZipFile("./test-data/docdb_xml_202402_CreateDelete_001.zip")
	ass.NoError(err)
}

func TestProcessDirectory(t *testing.T) {
	path := os.Getenv("DOCDB_BACKFILES_PATH")
	if len(path) == 0 {
		panic("no file path to the backfiles defined")
	}
	p := NewProcessor()
	p.IncludeAuthorities("EP")
	err := p.ProcessDirectory(path)
	if err != nil {
		t.Error(err)
	}
}
