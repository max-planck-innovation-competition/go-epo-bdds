package epo_docdb

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessBulkZipFile(t *testing.T) {
	// skipTest(t)
	log.SetLevel(log.TraceLevel)
	ass := assert.New(t)
	err := ProcessBulkZipFile(
		"./test-data/docdb_xml_202243_Amend_001.zip", "./test-data/xml")
	ass.NoError(err)
}
