package docdb

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type EpoDocDbResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deliveries  []struct {
		DeliveryID                  int         `json:"deliveryId"`
		DeliveryName                string      `json:"deliveryName"`
		DeliveryPublicationDatetime time.Time   `json:"deliveryPublicationDatetime"`
		DeliveryExpiryDatetime      interface{} `json:"deliveryExpiryDatetime"`
		Items                       []struct {
			ItemID                  int       `json:"itemId"`
			ItemName                string    `json:"itemName"`
			FileSize                string    `json:"fileSize"`
			FileChecksum            string    `json:"fileChecksum"`
			ItemPublicationDatetime time.Time `json:"itemPublicationDatetime"`
		} `json:"items"`
	} `json:"deliveries"`
}

// EpoDocDBProductEndpoint is the endpoint for the doc db product
var EpoDocDBProductEndpoint = "https://publication-bdds.apps.epo.org/bdds/bdds-bff-service/prod/api/products/3"

// GetDocDbFrontFileLinks returns the links to the front files of the doc db
func GetDocDbFrontFileLinks(tokenResponse TokenResponse) (response EpoDocDbResponse, err error) {
	// build token
	token := fmt.Sprintf("%s %s", tokenResponse.TokenType, tokenResponse.AccessToken)
	// create new http request with header and payload
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", EpoDocDBProductEndpoint, strings.NewReader(""))
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
