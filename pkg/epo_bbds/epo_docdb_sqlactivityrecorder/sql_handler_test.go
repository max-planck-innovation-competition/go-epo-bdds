package epo_docdb_sqlactivityrecorder

import (
	"fmt"
	"strconv"
	"testing"
)

// TODO
// Logger umbenennen
// Exchange files hinzufügen
// Löschen nach unten
// Einfügen in Handler
func TestSQLHandlerBreakMiddle(t *testing.T) {
	sqllogger := NewSqlActivityRecorder("log.db", "./", "C:\\docdb")
	fmt.Println(sqllogger)
	if sqllogger.IsDirectoryFinished() {
		return
	}

	for bulkindex := 1; bulkindex < 5; bulkindex++ {
		bulkname := "test" + strconv.Itoa(bulkindex)
		bulkstatus, _ := sqllogger.RegisterOrSkipZipFile(bulkname + ".zip")
		if bulkstatus == Done {
			continue
		}

		for xmlindex := 1; xmlindex < 5; xmlindex++ {

			xmlname := "frontfile" + strconv.Itoa(xmlindex)
			xmlstatus, _ := sqllogger.RegisterOrSkipXMLFile(bulkname+"_"+xmlname+".xml", "/DOC/files/")
			if bulkindex == 2 && xmlindex == 4 {
				return
			}
			if xmlstatus == Done {
				continue
			}

			sqllogger.MarkXMLAsFinished()
		}

		sqllogger.MarkZipFileAsFinished()
	}

	sqllogger.MarkProcessingDirectoryAsFinished()
}

func TestSQLHandlerFull(t *testing.T) {
	sqllogger := NewSqlActivityRecorder("log.db", "./", "C:\\docdb")
	sqllogger.SetSafeDelete(true)
	fmt.Println(sqllogger)
	if sqllogger.IsDirectoryFinished() {
		return
	}

	for bulkindex := 1; bulkindex < 5; bulkindex++ {
		bulkname := "test" + strconv.Itoa(bulkindex)
		bulkstatus, _ := sqllogger.RegisterOrSkipZipFile(bulkname + ".zip")
		if bulkstatus == Done {
			continue
		}

		for xmlindex := 1; xmlindex < 5; xmlindex++ {
			xmlname := "frontfile" + strconv.Itoa(xmlindex)
			xmlstatus, _ := sqllogger.RegisterOrSkipXMLFile(bulkname+"_"+xmlname+".xml", "/DOC/files/")
			if xmlstatus == Done {
				continue
			}

			for exchangeindex := 1; exchangeindex < 20; exchangeindex++ {
				exchangename := "exchangefile" + strconv.Itoa(xmlindex) + "-" + strconv.Itoa(exchangeindex)
				exchangestatus, _ := sqllogger.RegisterOrSkipExchangeLine(exchangename, exchangeindex)
				if exchangestatus == Done {
					continue
				}

				sqllogger.MarkExchangeAsFinished()
			}

			sqllogger.MarkXMLAsFinished()
		}

		sqllogger.MarkZipFileAsFinished()
	}

	sqllogger.MarkProcessingDirectoryAsFinished()
}
