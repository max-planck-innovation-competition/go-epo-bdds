package epo_docdb_sqllogger

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSQLHandler(t *testing.T) {
	sqllogger := NewSqlLogger("log.db", "./", "C:\\docdb")
	fmt.Println(sqllogger)
	//assert(sqllogger.GetDirectoryProcessStatus == NotProcessed)

	for bulkindex := 1; bulkindex < 5; bulkindex++ {
		bulkstatus, _ := sqllogger.RegisterOrSkipBulkFile("test" + strconv.Itoa(bulkindex) + ".zip")
		if bulkstatus == Done {
			break
		}

		for xmlindex := 1; xmlindex < 5; xmlindex++ {
			xmlstatus, _ := sqllogger.RegisterOrSkipBulkFile("frontfile" + strconv.Itoa(xmlindex) + ".xml")
			if xmlstatus == Done {
				break
			}
			sqllogger.MarkXMLAsFinished()
		}

		sqllogger.MarkBulkFileAsFinished()
	}

	sqllogger.MarkProcessingDirectoryAsFinished()
}
