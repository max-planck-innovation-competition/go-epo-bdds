package epo_bbds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// ErrNo200StatusCode is returned if the response status code is not 200
var ErrNo200StatusCode = errors.New("no 200 status code")

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
func GetAuthorizationToken() (token string, err error) {

	epoUserName := os.Getenv("EPO_USERNAME")
	if epoUserName == "" {
		err = errors.New("no epo username set")
		slog.With("err", err).Error("no epo username set")
		return
	}

	epoPassword := os.Getenv("EPO_PASSWORD")
	if epoPassword == "" {
		err = errors.New("no epo password set")
		slog.With("err", err).Error("no epo password set")
		return
	}

	payload := fmt.Sprintf("grant_type=password&username=%s&password=%s&scope=openid", epoUserName, epoPassword)

	// create new http request with header and payload
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", EpoLoginEndpoint, strings.NewReader(payload))
	if err != nil {
		slog.With("err", err).Error("failed to create new request")
		return
	}
	// add header
	req.Header.Set("Authorization", "Basic MG9hM3VwZG43YW41cE1JOE80MTc=")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.With("err", err).Error("failed to send request")
		return
	}
	// check status code
	if resp.StatusCode != 200 {
		err = ErrNo200StatusCode
		slog.With("err", err).With("statusCode", resp.StatusCode).Error("no 200 status code")
		return
	}
	// close response body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.With("err", err).Error("failed to close body")
		}
	}(resp.Body)
	var response TokenResponse
	// parse response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		slog.With("err", err).Error("failed to decode response")
		return
	}
	// check if token is contained in response
	if response.AccessToken == "" {
		err = ErrNoAccessToken
		slog.With("err", err).Error("no access token in response")
		return
	}
	return buildAuthToken(response)
}

// buildAuthToken builds the authorization token for the EPO API
func buildAuthToken(tokenResponse TokenResponse) (token string, err error) {
	// check if not empty
	if tokenResponse.TokenType == "" {
		err = ErrNoAccessToken
		slog.With("err", err).Error("no access token in response")
		return
	}
	// check if not empty
	if tokenResponse.AccessToken == "" {
		err = ErrNoAccessToken
		slog.With("err", err).Error("no access token in response")
		return
	}
	// build token
	token = fmt.Sprintf("%s %s", tokenResponse.TokenType, tokenResponse.AccessToken)
	return
}
