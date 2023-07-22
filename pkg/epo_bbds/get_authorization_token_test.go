package epo_bbds

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAuthorization(t *testing.T) {
	ass := assert.New(t)
	token, err := GetAuthorizationToken()
	ass.NoError(err)
	ass.NotEmpty(token)
	fmt.Println(token)
}
