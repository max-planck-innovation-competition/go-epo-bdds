package epo_bbds

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		"./test",
		resFrontFiles.Deliveries[0].Files[0].FileName,
	)
	if err != nil {
		t.Error(err)
	}
	ass.NoError(err)
}