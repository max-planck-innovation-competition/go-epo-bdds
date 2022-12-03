package bbds_api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDocDbFrontFileLinks(t *testing.T) {
	ass := assert.New(t)
	resToken, err := GetAuthorizationToken()
	ass.NoError(err)
	resFrontFiles, err := GetEpoBddsFileItems(resToken)
	ass.NoError(err)
	fmt.Println(resFrontFiles)
}
