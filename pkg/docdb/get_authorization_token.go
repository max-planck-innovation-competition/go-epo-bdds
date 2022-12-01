package docdb

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

// EpoLoginEndpoint is the endpoint for the EPO login
const EpoLoginEndpoint = "https://login.epo.org/oauth2/aus3up3nz0N133c0V417/v1/token"

// TokenResponse is the response from the EPO login endpoint
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenID     string `json:"id_token"`
}

// ErrNoAccessToken is returned if the token response does not contain an access token
var ErrNoAccessToken = errors.New("no access token")

// GetAuthorizationToken returns the authorization token for the EPO API
func GetAuthorizationToken() (response TokenResponse, err error) {
	payload := fmt.Sprintf("grant_type=password&username=%s&password=%s&scope=openid", os.Getenv("EPO_USERNAME"), os.Getenv("EPO_PASSWORD"))

	// create new http request with header and payload
	req, err := http.NewRequest("POST", EpoLoginEndpoint, strings.NewReader(payload))
	if err != nil {
		log.WithError(err).Error("failed to create new request")
		return
	}
	// add header
	req.Header.Set("Authorization", "Basic MG9hM3VwZG43YW41cE1JOE80MTc=")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err).Error("failed to send request")
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
	// check if token is contained in response
	if response.AccessToken == "" {
		err = ErrNoAccessToken
		log.WithError(err).Error("no access token in response")
		return
	}
	return
}
