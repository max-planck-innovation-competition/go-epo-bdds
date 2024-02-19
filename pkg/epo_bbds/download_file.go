package epo_bbds

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ErrCanNotDownload is thrown if the download is not possible
var ErrCanNotDownload = errors.New("can not download file")

// DownloadFile downloads a file from the bulk data service
func DownloadFile(token string, productID EpoBddsBProductID, deliveryID, fileID int, destinationFilePath, destinationFileName string) (err error) {
	// build endpoint url
	endpoint := fmt.Sprintf(EpoBddsFileEndpoint, string(productID), deliveryID, fileID)
	// create path if not exists
	err = os.MkdirAll(destinationFilePath, os.ModePerm)
	if err != nil {
		slog.With("err", err).Error("failed to create file path")
		return
	}
	// join file and filepath
	path := filepath.Join(destinationFilePath, destinationFileName)
	// create file
	out, err := os.Create(path)
	if err != nil {
		slog.With("err", err).Error("failed to create file")
		return
	}
	defer out.Close()
	// download file
	req, err := http.NewRequestWithContext(context.TODO(), "GET", endpoint, strings.NewReader(""))
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
	defer resp.Body.Close()
	// check status code
	if resp.StatusCode != http.StatusOK {
		slog.With("status", resp.Status).Error("failed to download file")
		err = ErrCanNotDownload
		slog.With("err", err).Error("failed to download file")
		return
	}

	// copy file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		slog.With("err", err).Error("failed to copy file")
		return
	}
	return
}

// DownloadAllFiles downloads all files from the bulk data service of a product
func DownloadAllFiles(token string, productID EpoBddsBProductID, destinationPath string) (err error) {
	// get back files
	deliveries, err := GetEpoBddsFileItems(token, productID)
	if err != nil {
		slog.With("err", err).Error("could not get files")
		return
	}
	for i, d := range deliveries.Deliveries {
		amountFiles := len(d.Files)
		for j, f := range d.Files {
			// check if file exists
			if pathExists(filepath.Join(destinationPath, f.FileName)) {
				slog.With("file", f.FileName, "no", j, "total", amountFiles).Info("file exists already")
				continue
			}
			slog.With("file", f.FileName, "no", j, "total", amountFiles).Info("start downloading")
			errDownload := DownloadFile(
				token,
				EpoDocDBFrontFilesProductID,
				deliveries.Deliveries[i].DeliveryID,
				deliveries.Deliveries[i].Files[j].FileID,
				destinationPath,
				deliveries.Deliveries[i].Files[j].FileName,
			)
			if errDownload != nil {
				slog.With("err", errDownload).Error("could not download file")
				err = errDownload
				return
			}
		}
	}
	slog.Info("Done")
	return
}

// pathExists checks if a path exists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
