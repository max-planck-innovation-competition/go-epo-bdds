package epo_bbds

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
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
		slog.With("err", err).Error("failed to read file")
		return
	}
	err = xml.Unmarshal(data, &indexObject)
	if err != nil {
		slog.With("err", err).Error("failed to unmarshal xml")
		return
	}
	return indexObject, nil
}

func readCsv(filePath string) ([][]string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file, %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.With("err", err).Error("couldn't close file")
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

/*
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
*/
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

type ExchangeDocuments struct {
	XMLName           xml.Name           `xml:"exchange-documents"`
	ExchangeDocuments []ExchangeDocument `xml:"exchange-document"`
}

type ExchangeDocument struct {
	XMLName xml.Name `xml:"exchange-document"`
}

// GenerateReplacerEntityMap transforms a given DTD file into a Replacer Entity map
func GenerateReplacerEntityMap(dtdFilePath string) (*strings.Replacer, error) {
	// Read the DTD file
	file, err := os.Open(dtdFilePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	entityMap := make(map[string]string)

	entityPattern := regexp.MustCompile(`<!ENTITY\s+(\w+)\s+"([^"]+)"`)

	// read and process each line of the DTD file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := entityPattern.FindStringSubmatch(line)
		if len(match) == 3 {
			entityName := match[1]
			entityValue := match[2]
			entityMap[entityName] = entityValue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// generate strings.Replacer instance using the entity map
	replacerArgs := make([]string, 0, len(entityMap)*2)
	for key, value := range entityMap {
		replacerArgs = append(replacerArgs, "&"+key+";", value)
	}

	return strings.NewReplacer(replacerArgs...), nil
}

// CountExchangeDocuments counts the exchange documents of an XML file
func CountExchangeDocuments(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	//dirname := "docdb_xml_202331_Amend_001"

	dtdname := "test-data/docdb-entities.dtd"
	replacer, _ := GenerateReplacerEntityMap(dtdname)

	preprocessedData := replacer.Replace(string(data))

	var exchangeDocs ExchangeDocuments
	err = xml.Unmarshal([]byte(preprocessedData), &exchangeDocs)
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

// CountZIPs unpacks zip file and counts DOCDB zip files within
func CountZIPs() {
	//
	zipFilePath := "test-data/docdb_xml_202331_Amend_001.zip"

	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(zipFile *zip.ReadCloser) {
		err := zipFile.Close()
		if err != nil {

		}
	}(zipFile)

	firstDirName := "docdb_xml_202331_Amend_001"
	rootDirName := "Root"
	docDirName := "DOC"

	var firstDir *zip.File
	for _, file := range zipFile.File {
		if file.Name == firstDirName+"/" {
			firstDir = file
			break
		}
	}

	if firstDir == nil {
		log.Fatal("First directory not found in the zip file.")
	}

	var rootDir *zip.File
	for _, file := range zipFile.File {
		if strings.HasPrefix(file.Name, firstDirName+"/") {
			if strings.TrimPrefix(file.Name, firstDirName+"/") == rootDirName+"/" {
				rootDir = file
				break
			}
		}
	}

	if rootDir == nil {
		log.Fatal("Root directory not found in the zip file.")
	}

	var docDir *zip.File
	for _, file := range zipFile.File {
		if strings.HasPrefix(file.Name, firstDirName+"/"+rootDirName+"/") {
			if strings.TrimPrefix(file.Name, firstDirName+"/"+rootDirName+"/") == docDirName+"/" {
				docDir = file
				break
			}
		}
	}

	if docDir == nil {
		log.Fatal("DOC directory not found inside the Root directory.")
	}

	count := 0
	for _, file := range zipFile.File {
		if strings.HasPrefix(file.Name, firstDirName+"/"+rootDirName+"/"+docDirName+"/") {
			if strings.HasPrefix(file.Name, firstDirName+"/"+rootDirName+"/"+docDirName+"/DOCDB") {
				count++
			}
		}
	}

	fmt.Printf("Number of zip files in %s directory starting with \"DOCDB\": %d\n", docDirName, count)
}

// Unzip unpacks a zip file
func Unzip(zipPath string, outputPath string) {

	dst := outputPath
	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(err)
	}
	defer func(archive *zip.ReadCloser) {
		err := archive.Close()
		if err != nil {
		}
	}(archive)

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		err = dstFile.Close()
		if err != nil {
			return
		}
		err = fileInArchive.Close()
		if err != nil {
			return
		}
	}
}

// CountExchDocs unpacks each zip file within a directory and counts its exchange documents
func CountExchDocs() {

	//zipPath := "test-data/docdb_xml_202331_Amend_001.zip"
	outputPath := "test-data/output"

	//Unzip(zipPath, outputPath)

	dirPath := outputPath + "/docdb_xml_202331_Amend_001/Root/DOC"

	dir, err := os.Open(dirPath)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return
	}
	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {
		}
	}(dir)

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "DOCDB") {
			dst := "test-data/output"

			//Unzip(dirPath + "/" + file.Name(), dst)
			count, err := CountExchangeDocuments(dst + "/" + strings.TrimSuffix(file.Name(), ".zip") + ".xml")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("File: %s, Exchange Documents Count: %d\n", file.Name(), count)
		}
	}
}

// CountExchDocsTemp unpacks each zip file within a temporary directory and counts its exchange documents
func CountExchDocsTemp() {

	zipPath := "test-data/docdb_xml_202331_Amend_001.zip"
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(zipFile *zip.ReadCloser) {
		err := zipFile.Close()
		if err != nil {

		}
	}(zipFile)

	tempDir, err := os.MkdirTemp("", "temp-docs")
	if err != nil {
		log.Fatal(err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {

		}
	}(tempDir)

	for _, file := range zipFile.File {
		if strings.HasPrefix(file.Name, "docdb_xml_202331_Amend_001/Root/DOC/DOCDB") {
			fmt.Println(file.Name)
			dstPath := filepath.Join(tempDir, filepath.Base(file.Name))
			dstFile, err := os.Create(dstPath)
			if err != nil {
				log.Printf("Error creating destination file: %s", err)
				continue
			}
			srcFile, err := file.Open()
			if err != nil {
				log.Printf("Error opening source file: %s", err)
				err := dstFile.Close()
				if err != nil {
					return
				}
				continue
			}
			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				log.Printf("Error copying file contents: %s", err)
			}
			err = dstFile.Close()
			if err != nil {
				return
			}
			err = srcFile.Close()
			if err != nil {
				return
			}

			xmlTempDir := filepath.Join(tempDir, "xmltemp")
			err = os.Mkdir(xmlTempDir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			Unzip(dstPath, xmlTempDir)

			xmlFiles, err := os.ReadDir(xmlTempDir)
			if err != nil {
				log.Printf("Error reading extracted XML directory: %s", err)
				continue
			}

			for _, xmlFile := range xmlFiles {
				xmlFilePath := filepath.Join(xmlTempDir, xmlFile.Name())

				count, err := CountExchangeDocuments(xmlFilePath)
				if err != nil {
					log.Printf("Error counting exchange documents: %s", err)
					continue
				}
				fmt.Printf("File: %s, Exchange Documents Count: %d\n", xmlFile.Name(), count)
			}

			// clean up
			err = os.RemoveAll(xmlTempDir)
			if err != nil {
				log.Printf("Error cleaning up extracted XML directory: %s", err)
			}
		}
	}
}
