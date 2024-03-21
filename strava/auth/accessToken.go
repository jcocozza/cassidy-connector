package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	AuthorizationUrl       string = "https://www.strava.com/oauth/token"
	GrantAuthorizationCode string = "authorization_code"
	GrantRefreshToken      string = "refresh_token"
)
// This contains the user's short-lived access token which is used to access data.
// When it expires, use the refresh token to get a new access token
//
// This struct is obtained in 1 of 2 ways.
// (1) by possessing an existing refresh token and getting a new access token.
// (2) via user authorization, whereby an auth code is issued which is used to get the access token.
type GetAccessTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
// Create the payload for a refresh token
func makePayloadRefreshToken(refreshToken string) ([]byte, error) {
	payload := map[string]string{
		"client_id":     ClientId,
		"client_secret": ClientSecret,
		"refresh_token": refreshToken,
		"grant_type":    GrantRefreshToken,
		"f":             "json",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return payloadBytes, nil
}
// Get the access token from strava for an existing application
//
// The user must already have granted access to the application to use this.
// (Otherwise, must authenticate first)
func RefreshAccessToken(refreshToken string) (*GetAccessTokenResponse, error) {
	// Create the data
	payload, err := makePayloadRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	// Create a request
	req, err1 := http.NewRequest("POST", AuthorizationUrl, bytes.NewBuffer(payload))
	if err1 != nil {
		return nil, err1
	}
	req.Header.Set("Content-Type", "application/json")
	// Make the request
	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return nil, err2
	}
	defer resp.Body.Close()
	// Read response body
	body, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		return nil, err3
	}
	// Load response into struct
	var response GetAccessTokenResponse
	err4 := json.Unmarshal(body, &response)
	if err4 != nil {
		return nil, err4
	}
	return &response, nil
}