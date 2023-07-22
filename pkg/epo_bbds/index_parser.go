package epo_bbds

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// index struct

type DocdbPackageIndex struct {
	XMLName          xml.Name `xml:"docdb-package-index"`
	Text             string   `xml:",chardata"`
	ID               string   `xml:"id,attr"`
	DateProduced     string   `xml:"date-produced,attr"`
	DtdVersion       string   `xml:"dtd-version,attr"`
	File             string   `xml:"file,attr"`
	ProducedBy       string   `xml:"produced-by,attr"`
	VolumeID         string   `xml:"volume-id,attr"`
	DocdbPackageFile []struct {
		Text         string `xml:",chardata"`
		ID           string `xml:"id,attr"`
		Format       string `xml:"format,attr"`
		Size         string `xml:"size,attr"`
		Filename     string `xml:"filename"`
		FileLocation struct {
			Text     string `xml:",chardata"`
			Relative string `xml:"relative,attr"`
		} `xml:"file-location"`
		DocdbDocRange struct {
			Text                 string `xml:",chardata"`
			Country              string `xml:"country,attr"`
			DocdbFirstDocInRange struct {
				Text      string `xml:",chardata"`
				Country   string `xml:"country"`
				DocNumber string `xml:"doc-number"`
				Kind      string `xml:"kind"`
				Date      string `xml:"date"`
			} `xml:"docdb-first-doc-in-range"`
			DocdbLastDocInRange struct {
				Text      string `xml:",chardata"`
				Country   string `xml:"country"`
				DocNumber string `xml:"doc-number"`
				Kind      string `xml:"kind"`
				Date      string `xml:"date"`
			} `xml:"docdb-last-doc-in-range"`
		} `xml:"docdb-doc-range"`
	} `xml:"docdb-package-file"`
}

func ParseIndexXML(filename string) (indexObject DocdbPackageIndex, err error) {
	// Read the XML file
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Failed to read file: %s\n", err.Error())
		return
	}
	err = xml.Unmarshal(data, &indexObject)
	if err != nil {
		log.Printf("Failed to parse XML: %s\n", err.Error())
		return
	}
	return indexObject, nil
}

type ExchangeDocuments struct {
	XMLName           xml.Name           `xml:"exchange-documents"`
	ExchangeDocuments []ExchangeDocument `xml:"exchange-document"`
}

type ExchangeDocument struct {
	XMLName xml.Name `xml:"exchange-document"`
}

func CountExchangeDocuments(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	var exchangeDocs ExchangeDocuments
	err = xml.Unmarshal(data, &exchangeDocs)
	if err != nil {
		return 0, err
	}

	count := len(exchangeDocs.ExchangeDocuments)
	return count, nil
}

/*
statistics csv:
DOCA095,202324,DOCDB-202324,616732,478
Status,NrOfPNStatus,CC,KC,FirstPN,LastPN,FirstPubDate,LastPubDate,NrOfPN
A,616732,AP,A,36,4072,19881206,20170316,24
A,616732,AP,A0,8500014,2017009807,19850801,20170331,63
A,616732,AR,A1,000008,248143,19730919,20230419,863
A,616732,AR,A2,023762,240811,19910228,20230426,90
A,616732,AR,A4,067548,123822,20091014,20230118,2
A,616732,AT,A,A4386,A13742003,19761215,20060115,32
...

*/

type StatsStruct struct {
	Status       string
	NrOfPNStatus int
	CC           string
	KC           string
	FirstPN      int
	LastPN       int
	FirstPubDate time.Time
	LastPubDate  time.Time
	NrOfPN       int
}

func readCsv(filePath string) ([][]string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file, %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	scanner := bufio.NewScanner(file)

	//two dimensional slice
	var data [][]string

	for scanner.Scan() {
		row := scanner.Text()
		fields := strings.Split(row, ",")

		data = append(data, fields)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("couldn't read csv file: %s", err)
	}

	return data, nil
}

func deleteFirstLine(inputFile, outputFile string) error {
	//opens file as txt, deletes first line, returns mod output file

	input, err := os.Open(inputFile)

	if err != nil {
		return fmt.Errorf("couldn't open file, %s", err)
	}
	defer func(input *os.File) {
		err := input.Close()
		if err != nil {
		}
	}(input)

	output, err := os.Create(outputFile)

	if err != nil {
		return fmt.Errorf("couldn't create output file, %s", err)
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
		}
	}(output)

	scanner := bufio.NewScanner(input)

	if scanner.Scan() {
		for scanner.Scan() {
			line := scanner.Text()
			_, err := fmt.Fprintln(output, line)
			if err != nil {
				return fmt.Errorf("couldn't write to output file, %s", err)
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("couldn't read input file, %s", err)
		}
	}
	return nil
}

func readCsvToStruct(filePath string) ([]StatsStruct, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file, %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	// bufio bc csv reader can't handle the first line being different
	scanner := bufio.NewScanner(file)

	//skip first line (different amount of columns than rest)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("couldn't read the first line, %s", err)
		}
	}

	var data []StatsStruct

	for scanner.Scan() {
		row := scanner.Text()
		fields := strings.Split(row, ",")
		//string conversion
		nrOfPNStatus, _ := strconv.Atoi(fields[1])
		firstPN, _ := strconv.Atoi(fields[4])
		lastPN, _ := strconv.Atoi(fields[5])
		nrOfPN, _ := strconv.Atoi(fields[8])
		firstPubDate, _ := time.Parse("20060102", fields[6])
		lastPubDate, _ := time.Parse("20060102", fields[7])

		csvRow := StatsStruct{
			Status:       fields[0],
			NrOfPNStatus: nrOfPNStatus,
			CC:           fields[2],
			KC:           fields[3],
			FirstPN:      firstPN,
			LastPN:       lastPN,
			FirstPubDate: firstPubDate,
			LastPubDate:  lastPubDate,
			NrOfPN:       nrOfPN,
		}

		data = append(data, csvRow)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("couldn't read CSV file, %s", err)
	}

	return data, nil
}
