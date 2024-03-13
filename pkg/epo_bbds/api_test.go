package epo_bbds

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
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

	// get front files by id
	t.Log("Getting EPO File Items")
	resFrontFiles, err := GetEpoBddsFileItems(token, EpoDocDBFrontFilesProductID)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Got File Items")
	}

	res2B, _ := json.Marshal(resFrontFiles)
	os.WriteFile("big_marhsall.json", res2B, os.ModePerm)

	fmt.Println(string(res2B))
	ass.NoError(err)
}
