package state_handler

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSQLHandlerBreakMiddle(t *testing.T) {
	// TODO: windows directory as env variable
	stateHandler := NewStateHandler("log.db", "./", "C:\\docdb")
	fmt.Println(stateHandler)
	if stateHandler.IsDirectoryFinished() {
		return
	}

	for bulkIndex := 1; bulkIndex < 5; bulkIndex++ {
		bulkName := "test" + strconv.Itoa(bulkIndex)
		bulkStatus, _ := stateHandler.RegisterOrSkipZipFile(bulkName + ".zip")
		if bulkStatus == Done {
			continue
		}

		for xmlIndex := 1; xmlIndex < 5; xmlIndex++ {

			xmlName := "frontfile" + strconv.Itoa(xmlIndex)
			xmlStatus, _ := stateHandler.RegisterOrSkipXMLFile(bulkName+"_"+xmlName+".xml", "/DOC/files/")
			if bulkIndex == 2 && xmlIndex == 4 {
				return
			}
			if xmlStatus == Done {
				continue
			}

			stateHandler.MarkXMLAsFinished()
		}

		stateHandler.MarkZipFileAsFinished()
	}

	stateHandler.MarkProcessingDirectoryAsFinished()
}

func TestSQLHandlerFull(t *testing.T) {
	sqllogger := NewStateHandler("log.db", "./", "C:\\docdb")
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
