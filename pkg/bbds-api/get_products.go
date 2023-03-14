package bbds_api

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

// EpoProductsEndpoint is the endpoint for the products
var EpoProductsEndpoint = "https://publication-bdds.apps.epo.org/bdds/bdds-bff-service/prod/api/products/"

// EpoProductItem is the item of the products response
type EpoProductItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetProducts returns the products
func GetProducts(token string) (response []EpoProductItem, err error) {
	// create new http request with header and payload
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", EpoProductsEndpoint, strings.NewReader(""))
	if err != nil {
		log.WithError(err).Error("failed to create new request")
		return
	}
	// add header
	req.Header.Set("Authorization", token)

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err).Error("failed to send request")
		return
	}
	// check status code
	if resp.StatusCode != 200 {
		err = ErrNo200StatusCode
		log.WithError(err).Error("server responded with non 200 status code")
		return
	}
	// close response body
	defer resp.Body.Close()
	// parse response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.WithError(err).Error("failed to parse response")
		return
	}
	return
}
