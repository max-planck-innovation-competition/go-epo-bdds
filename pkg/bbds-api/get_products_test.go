package bbds_api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetProducts(t *testing.T) {
	ass := assert.New(t)
	resToken, err := GetAuthorizationToken()
	ass.NoError(err)

	resProducts, err := GetProducts(resToken)
	ass.NoError(err)
	fmt.Println(resProducts)
}
