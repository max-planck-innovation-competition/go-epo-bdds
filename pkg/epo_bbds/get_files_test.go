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
}
