package epo_bbds

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestCountExchangeDocuments(t *testing.T) {
	ass := assert.New(t)
	// Call the ParseXMLFile function with the file path
	response, err := CountExchangeDocuments("test-data/DOCDB-202324-Amend-PubDate20230609AndBefore-AP-0001.xml")

	fmt.Printf("Number of Exchange Documents in Bulkfile: %d\n", response)
	ass.NoError(err)
	if err != nil {
		t.Error(err)
	}
}

func TestParseIndexXML(t *testing.T) {
	ass := assert.New(t)

	response, err := ParseIndexXML("test-data/index.xml")

	ass.NoError(err)
	if err != nil {
		t.Error(err)
	}
	// ass.Equal("document-id", response.DocdbPackageFile)

	fmt.Printf("ID: %s\n", response.ID)
	fmt.Printf("Name: %s\n", response.DocdbPackageFile)

}

func TestReadCsv(t *testing.T) {
	// works with this
	filePath1 := "test-data/test.csv"

	// actual stat file, this only prints one line?
	filePath2 := "test-data/statistics_202324_Amend_001.csv"

	csvData1, err := readCsv(filePath1)
	if err != nil {
		t.Fatalf("Error reading %s: %s", filePath1, err)
	}

	//line by line
	fmt.Printf("data from %s\n", filePath1)
	fmt.Printf("Number of lines: %d\n", len(csvData1))
	/*
		for _, row := range csvData1 {
			fmt.Println(strings.Join(row, ","))
		}
	*/

	csvData2, err := readCsv(filePath2)
	if err != nil {
		t.Fatalf("Error reading %s: %s", filePath2, err)
	}

	//line by line
	fmt.Printf("data from %s\n", filePath2)
	// Count the number of lines in csvData

	fmt.Printf("Number of lines: %d\n", len(csvData2))
	for _, row := range csvData2 {
		fmt.Println(strings.Join(row, ","))
	}

	if !reflect.DeepEqual(csvData1, csvData2) {
		t.Errorf("Data extracted from %s and %s did not match even though it should", filePath1, filePath2)
	}
}

func TestReadCsvToStruct(t *testing.T) {

	filePath1 := "test-data/test.csv"
	// actual stat file:
	filePath2 := "test-data/statistics_202324_Amend_001.csv"

	statistics1, err := readCsvToStruct(filePath1)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, row := range statistics1 {
		if i > 0 {
			fmt.Printf("%+v\n", row)
		}
	}

	statistics2, err := readCsvToStruct(filePath2)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, row := range statistics2 {
		if i > 0 {
			fmt.Printf("%+v\n", row)
		}
	}

	if !reflect.DeepEqual(statistics1, statistics2) {
		t.Errorf("Data extracted from %s and %s did not match even though it should", filePath1, filePath2)
	}

	/*
		//specific column
		for i, row := range statistics {
			if i > 0 {
				fmt.Println(row.NrOfPN)
			}
		}
	*/

}
