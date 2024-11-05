package epo_bbds

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestDownloadDocDbFrontFileWithEncodingIssues(t *testing.T) {
	ass := assert.New(t)

	t.Log("Getting EPO Token")
	// get token
	token, err := GetAuthorizationToken()
	ass.NoError(err)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Success")
	}

	// get front files
	t.Log("Getting EPO File Items")
	resFrontFiles, err := GetEpoBddsFileItems(token, EpoDocDBFrontFilesProductID)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Got File Items")
	}
	ass.NoError(err)

	// download front files
	t.Log("Downloading Files")
	err = DownloadFile(token,
		EpoDocDBFrontFilesProductID,
		resFrontFiles.Deliveries[0].DeliveryID,
		resFrontFiles.Deliveries[0].Files[0].FileID,
		"./test-data",
		resFrontFiles.Deliveries[0].Files[0].FileName,
	)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Files Downloaded")
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
	destinationPath := os.Getenv("DOCDB_BACKFILES_PATH")
	if len(destinationPath) == 0 {
		t.Error("no file path found")
		return
	}
	_, err := DownloadAllFiles(EpoDocDBBackFilesProductID, destinationPath)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDownloadDocDbFrontFiles(t *testing.T) {
	destinationPath := os.Getenv("DOCDB_FRONTFILES_PATH")
	if len(destinationPath) == 0 {
		t.Error("no file path found")
		return
	}
	_, err := DownloadAllFiles(EpoDocDBFrontFilesProductID, destinationPath)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDownloadPatstat(t *testing.T) {
	destinationPath := os.Getenv("PATSTAT_FILES_PATH")
	if len(destinationPath) == 0 {
		t.Error("no file path found")
		return
	}
	_, err := DownloadAllFiles(EpoPatstatGlobalProductID, destinationPath)
	if err != nil {
		t.Error(err)
		return
	}
}
