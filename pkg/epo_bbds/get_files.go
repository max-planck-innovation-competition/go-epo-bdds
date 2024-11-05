package epo_bbds

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// EpoDocDbFileItem is a single item of the doc db
type EpoDocDbFileItem struct {
	FileID                  int       `json:"fileId"`
	FileName                string    `json:"fileName"`
	FileSize                string    `json:"fileSize"`
	FileChecksum            string    `json:"fileChecksum"`
	ItemPublicationDatetime time.Time `json:"itemPublicationDatetime"`
}
type EpoDocDbFileItems []EpoDocDbFileItem

func (a EpoDocDbFileItems) Len() int      { return len(a) }
func (a EpoDocDbFileItems) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a EpoDocDbFileItems) Less(i, j int) bool {
	return a[i].FileID < a[j].FileID
}

// EpoProductDeliveriesResponse is the response from the epo doc db
type EpoProductDeliveriesResponse struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Deliveries  []EpoProductDelivery `json:"deliveries"`
}

// EpoProductDelivery is a single delivery of the epo bbds
type EpoProductDelivery struct {
	DeliveryID                  int                `json:"deliveryId"`
	DeliveryName                string             `json:"deliveryName"`
	DeliveryPublicationDatetime time.Time          `json:"deliveryPublicationDatetime"`
	DeliveryExpiryDatetime      *time.Time         `json:"deliveryExpiryDatetime"`
	Files                       []EpoDocDbFileItem `json:"files"`
}

// GetEpoBddsFileItems returns the links to the front files of the doc db
func GetEpoBddsFileItems(token string, productID EpoBddsBProductID) (response EpoProductDeliveriesResponse, err error) {

	// build endpoint url
	endpoint := fmt.Sprintf(EpoBddsProductEndpoint, string(productID))

	// create new http request with header and payload
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, strings.NewReader(""))
	if err != nil {
		slog.With("err", err).Error("failed to create new request")
		return
	}
	// add header
	req.Header.Set(AuthHeader, token)

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.With("err", err).Error("failed to send request")
		return
	}
	// check status code
	if resp.StatusCode != 200 {
		err = ErrNo200StatusCode
		slog.With("err", err).With("statusCode", resp.StatusCode).Error("server responded with non 200 status code")
		return
	}
	// close response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.With("err", err).Error("failed to close body")
		}
	}(resp.Body)
	// parse response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		slog.With("err", err).Error("failed to parse response")
		return
	}

	return
}
