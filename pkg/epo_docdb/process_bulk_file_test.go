package epo_docdb

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessBulkZipFile(t *testing.T) {
	// skipTest(t)
	log.SetLevel(log.TraceLevel)
	p := NewProcessor()
	err := p.ProcessBulkZipFile("./test-data/docdb_xml_202402_CreateDelete_001.zip")
	if err != nil {
		t.Error(err)
	}
}

func TestProcessBulkZipFile2023(t *testing.T) {
	// skipTest(t)
	log.SetLevel(log.TraceLevel)
	ass := assert.New(t)
	p := NewFileExportProcessor("./test-data/xml")
	err := p.ProcessBulkZipFile("./test-data/docdb_xml_202402_CreateDelete_001.zip")
	ass.NoError(err)
}

func TestProcessEpFiles2023(t *testing.T) {
	// skipTest(t)
	log.SetLevel(log.TraceLevel)
	ass := assert.New(t)
	p := NewFileExportProcessor("./test-data/xml/eps")
	p.IncludeAuthorities("EP")
	err := p.ProcessBulkZipFile("./test-data/docdb_xml_202402_CreateDelete_001.zip")
	ass.NoError(err)
}

func TestProcessDirectory(t *testing.T) {
	p := NewProcessor()
	p.IncludeAuthorities("EP")
	err := p.ProcessDirectory("/media/oem/Volume/docdb/backfiles")
	if err != nil {
		t.Error(err)
	}
}
