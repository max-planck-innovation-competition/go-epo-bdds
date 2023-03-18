package epo_bbds

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
	resFrontFiles, err := GetEpoBddsFileItems(resToken, EpoDocDBFrontFilesProductID)
	ass.NoError(err)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(resFrontFiles)
}

func TestGetDocDbBackFileLinks(t *testing.T) {
	ass := assert.New(t)
	resToken, err := GetAuthorizationToken()
	ass.NoError(err)
	if err != nil {
		t.Error(err)
		return
	}
	resFrontFiles, err := GetEpoBddsFileItems(resToken, EpoDocDBBackFilesProductID)
	ass.NoError(err)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(resFrontFiles)
	for _, d := range resFrontFiles.Deliveries {
		fmt.Println(d.DeliveryName)
		for _, f := range d.Files {
			fmt.Println("\t"+f.FileName, f.FileSize)
		}
	}
}
