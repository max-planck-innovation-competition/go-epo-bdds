package docdb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDocDbFrontFiles(t *testing.T) {
	ass := assert.New(t)
	resToken, err := GetAuthorizationToken()
	ass.NoError(err)

	resFrontFiles, err := GetDocDbFrontFileLinks(resToken)
	ass.NoError(err)
	fmt.Println(resFrontFiles)
}
