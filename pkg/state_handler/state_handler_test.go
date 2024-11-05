package state_handler

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var testDir = ""

const testDbName = "test.sqlite"

var testDatabasePath = ""

func init() {
	// get the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testDir = filepath.Join(pwd, "test-data")
	testDatabasePath = filepath.Join(testDir, testDbName)
}

func TestSQLHandlerBreakMiddle(t *testing.T) {
	// defer func(name string) {
	// 	err := os.Remove(testDatabasePath)
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// }(testDir)

	stateHandler := New(testDbName, testDir, testDir)
	if stateHandler.IsDirectoryFinished() {
		return
	}

	// generate synthetic data
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
	stateHandler := New(testDbName, testDir, testDir)
	stateHandler.SetSafeDelete(true)
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
			if xmlStatus == Done {
				continue
			}

			for exchangeIndex := 1; exchangeIndex < 20; exchangeIndex++ {
				exchangeName := "exchangefile" + strconv.Itoa(xmlIndex) + "-" + strconv.Itoa(exchangeIndex)
				exchangeStatus, _ := stateHandler.RegisterOrSkipExchangeLine(exchangeName, exchangeIndex)
				if exchangeStatus == Done {
					continue
				}

				stateHandler.MarkExchangeFileAsFinished()
			}

			stateHandler.MarkXMLAsFinished()
		}

		stateHandler.MarkZipFileAsFinished()
	}

	stateHandler.MarkProcessingDirectoryAsFinished()
}

func TestStateHandler(t *testing.T) {
	defer func(name string) {
		err := os.Remove(testDatabasePath)
		if err != nil {
			t.Error(err)
		}
	}(testDir)

	// create artificial data
	todos := []string{
		"a",
		"a/a1",
		"a/a1/a1.xml",
		"a/a2/a2.xml",
		"b",
		"b/b1",
		"b/b1/b1.xml",
		"b/b2/b2.xml",
		"c",
		"c/c1",
		"c/c1/c1.xml",
		"c/c2/c2.xml",
	}

	// create a new state handler
	stateHandler := New(testDbName, testDir, testDir)

	for _, todo := range todos {
		// register the file
		skip, err := stateHandler.RegisterOrSkip(todo)
		if err != nil {
			t.Error(err)
		}
		if skip == true {
			t.Error("should not be done")
		}
	}

	printState(stateHandler, t)

	// mark the file as finished
	err := stateHandler.MarkAsDone("a/a1/a1.xml")
	if err != nil {
		t.Error(err)
	}
	// check if the file is done
	skip, err := stateHandler.RegisterOrSkip("a/a1/a1.xml")
	if err != nil {
		t.Error(err)
	}
	if !skip {
		t.Error("should be done")
	}

	// mark the subdirectory as finished
	err = stateHandler.MarkAsDone("a/a1")
	if err != nil {
		t.Error(err)
	}

	// check if the file is done
	skip, err = stateHandler.RegisterOrSkip("a/a1/a1.xml")
	if err != nil {
		t.Error(err)
	}
	if skip == false {
		t.Error("should be done")
	}
	// check the sibling file
	skip, err = stateHandler.RegisterOrSkip("a/a2/a2.xml")
	if err != nil {
		t.Error(err)
	}
	if skip == true {
		t.Error("should not done")
	}
	//
	skip, err = stateHandler.RegisterOrSkip("b")
	if err != nil {
		t.Error(err)
	}
	if skip == true {
		t.Error("should not be done")
	}

	// mark the file as finished
	_ = stateHandler.MarkAsDone("a")
	_ = stateHandler.MarkAsDone("b")

	// check the sibling file
	skip, err = stateHandler.RegisterOrSkip("b/b1/b1.xml")
	if err != nil {
		t.Error(err)
	}
	if skip == false {
		t.Error("should be done")
	}

	// check the sibling file
	skip, err = stateHandler.RegisterOrSkip("b/b2/b2.xml")
	if err != nil {
		t.Error(err)
	}
	if skip == false {
		t.Error("should be done")
	}

	// check c
	skip, err = stateHandler.RegisterOrSkip("c")
	if err != nil {
		t.Error(err)
	}
	if skip == true {
		t.Error("should not be done")
	}

	// check lower c file
	skip, err = stateHandler.RegisterOrSkip("c/c1/c1.xml")
	if err != nil {
		t.Error(err)
	}
	if skip == true {
		t.Error("should not be done")
	}

	printState(stateHandler, t)

}

func printState(stateHandler *StateHandler, t *testing.T) {
	// get the full database
	all := []*Object{}
	err := stateHandler.db.Find(&all).Error
	if err != nil {
		t.Error(err)
	}
	for _, a := range all {
		t.Log(a)
	}
}
