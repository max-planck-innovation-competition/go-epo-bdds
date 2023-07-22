package epo_bbds

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strconv"
	"testing"
)

func TestGetDocDbFrontFileLinks(t *testing.T) {
	ass := assert.New(t)
	resToken, err := GetAuthorizationToken()
	ass.NoError(err)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := GetEpoBddsFileItems(resToken, EpoDocDBFrontFilesProductID)
	ass.NoError(err)
	if err != nil {
		t.Error(err)
		return
	}
	printFiles(res)
}

func TestGetDocDbBackFileLinks(t *testing.T) {
	ass := assert.New(t)
	resToken, err := GetAuthorizationToken()
	ass.NoError(err)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := GetEpoBddsFileItems(resToken, EpoDocDBBackFilesProductID)
	ass.NoError(err)
	if err != nil {
		t.Error(err)
		return
	}
	printFiles(res)
}

func TestGetPatstatGlobalFileLinks(t *testing.T) {
	ass := assert.New(t)
	resToken, err := GetAuthorizationToken()
	ass.NoError(err)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := GetEpoBddsFileItems(resToken, EpoPatstatGlobalProductID)
	ass.NoError(err)
	if err != nil {
		t.Error(err)
		return
	}
	printFiles(res)
}

func printFiles(res EpoProductDeliveriesResponse) {
	totalFileSize := int64(0)
	for _, d := range res.Deliveries {
		fmt.Printf("%d \t %s\n", d.DeliveryID, d.DeliveryName)
		for _, f := range d.Files {
			size, err := parseFileSize(f.FileSize)
			if err != nil {
				fmt.Println(err)
				break
			} else {
				totalFileSize += size
			}
			fmt.Printf("\t %d \t %s %s \n", f.FileID, f.FileName, f.FileSize)
		}
	}
	fmt.Println("\n\nTotal file size in GB:", float64(totalFileSize)/1024/1024/1024)
}

func parseFileSize(size string) (result int64, err error) {
	// e.g.  1.7 GB
	// extract number and unit
	re := regexp.MustCompile(`(\d+\.?\d*)\s*(\w+)`)
	matches := re.FindStringSubmatch(size)
	if len(matches) != 3 {
		return 0, fmt.Errorf("could not parse size: %s", size)
	}
	// number
	number, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, err
	}
	// unit
	unit := matches[2]

	switch unit {
	case "GB":
		result = int64(number * 1024 * 1024 * 1024)
		return
	case "MB":
		result = int64(number * 1024 * 1024)
		return
	case "KB":
		result = int64(number * 1024)
	}
	return
}
