package epo_bbds

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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
	// join file and filepath
	path := filepath.Join(destinationFilePath, destinationFileName)
	// create file
	out, err := os.Create(path)
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
	req.Header.Set(AuthHeader, token)
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
