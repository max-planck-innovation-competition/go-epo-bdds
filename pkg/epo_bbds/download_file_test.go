package epo_bbds

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

func TestDownloadDocDbFrontFileWithEncodingIssues(t *testing.T) {
	ass := assert.New(t)

	// get token
	token, err := GetAuthorizationToken()
	ass.NoError(err)
	if err != nil {
		t.Error(err)
	}

	// get front files
	resFrontFiles, err := GetEpoBddsFileItems(token, EpoDocDBFrontFilesProductID)
	if err != nil {
		t.Error(err)
	}
	ass.NoError(err)

	// download front files
	err = DownloadFile(token,
		EpoDocDBFrontFilesProductID,
		resFrontFiles.Deliveries[0].DeliveryID,
		resFrontFiles.Deliveries[0].Files[0].FileID,
		"./test-data",
		resFrontFiles.Deliveries[0].Files[0].FileName,
	)
	if err != nil {
		t.Error(err)
	}
	ass.NoError(err)
}

func TestDownloadDocDbBackFile(t *testing.T) {
	ass := assert.New(t)

	// get token
	token, err := GetAuthorizationToken()
	ass.NoError(err)
	if err != nil {
		t.Error(err)
	}

	// get front files
	backFiles, err := GetEpoBddsFileItems(token, EpoDocDBBackFilesProductID)
	if err != nil {
		t.Error(err)
	}
	ass.NoError(err)

	backFileDelivery := EpoProductDelivery{}
	ok := false
	for i := range backFiles.Deliveries {
		if strings.Contains(backFiles.Deliveries[i].DeliveryName, "DOCDB Back file") {
			backFileDelivery = backFiles.Deliveries[i]
			ok = true
		}
	}
	if !ok {
		t.Failed()
		return
	}

	// sort files
	files := EpoDocDbFileItems(backFileDelivery.Files)
	fmt.Println(files)
	sort.Sort(files)
	fmt.Println(files)
	// download front files
	err = DownloadFile(token,
		EpoDocDBBackFilesProductID,
		backFileDelivery.DeliveryID,
		files[0].FileID,
		"./test-data",
		files[0].FileName,
	)
	if err != nil {
		t.Error(err)
	}
	ass.NoError(err)
}

func TestDownloadDocDbFrontFile(t *testing.T) {
	ass := assert.New(t)
	// get token
	token, err := GetAuthorizationToken()
	ass.NoError(err)
	if err != nil {
		t.Error(err)
	}

	// get front files
	resFrontFiles, err := GetEpoBddsFileItems(token, EpoDocDBFrontFilesProductID)
	if err != nil {
		t.Error(err)
	}
	ass.NoError(err)

	// download front files
	err = DownloadFile(token,
		EpoDocDBFrontFilesProductID,
		resFrontFiles.Deliveries[0].DeliveryID,
		resFrontFiles.Deliveries[0].Files[0].FileID,
		"./test-data",
		resFrontFiles.Deliveries[0].Files[0].FileName,
	)
	if err != nil {
		t.Error(err)
	}
	ass.NoError(err)
}

func TestDownloadDocDbBackFiles(t *testing.T) {
	destinationPath := "/media/oem/Volume/docdb/backfiles"
	err := DownloadAllFiles(EpoDocDBBackFilesProductID, destinationPath)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDownloadDocDbFrontFiles(t *testing.T) {
	destinationPath := "/media/seb/T18-1/docdb-frontfiles"
	err := DownloadAllFiles(EpoDocDBFrontFilesProductID, destinationPath)
	if err != nil {
		t.Error(err)
		return
	}
}
