package epo_bbds

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAuthorization(t *testing.T) {
	ass := assert.New(t)
	res, err := GetAuthorizationToken()
	ass.NoError(err)
	ass.NotEmpty(res)
}