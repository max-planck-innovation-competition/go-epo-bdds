package bbds_api

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
)

// EpoBddsFileEndpoint is the endpoint for the docdb frontfiles bucket
// GET https://publication-bdds.apps.epo.org/bdds/bdds-bff-service/prod/api/products/%s/delivery/%s/file/%s/download
var EpoBddsFileEndpoint = "https://publication-bdds.apps.epo.org/bdds/bdds-bff-service/prod/api/products/%s/delivery/%d/file/%d/download"

// DownloadFile downloads a file from the bulk data service
func DownloadFile(token string, productID EpoBddsBProductID, deliveryID, fileID int, destinationFilePath, destinationFileName string) (err error) {
	// build endpoint url
	endpoint := fmt.Sprintf(EpoBddsFileEndpoint, string(productID), deliveryID, fileID)
	// create path if not exists
	err = os.MkdirAll(destinationFilePath, os.ModePerm)
	if err != nil {
		log.WithError(err).Error("failed to create file path")
		return
	}
	// create file
	out, err := os.Create("./" + destinationFileName)
	if err != nil {
		log.WithError(err).Error("failed to create file")
		return
	}
	defer out.Close()
	// download file
	req, err := http.NewRequestWithContext(context.TODO(), "GET", endpoint, strings.NewReader(""))
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
	defer resp.Body.Close()
	// copy file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.WithError(err).Error("failed to copy file")
		return
	}
	return
}
