package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	//"github.com/jcocozza/cassidy-connector/strava/internal"
	"github.com/jcocozza/cassidy-connector/strava/internal/auth"
)

const (
	responseType string = "code"
	redirectUri string = "http://localhost/exchange_token"
	approvalPrompt string = "force"
	scope = "activity:read_all"
	approvalUrlFormat string = "https://www.strava.com/oauth/authorize?client_id=%s&response_type=%s&redirect_uri=%s&approval_prompt=%s&scope=%s"
)

type Authenticator interface {
	// Format the ApprovalUrlFormat
	GenerateApprovalUrl() string
	// This is for the FIRST TIME getting the access token.
	//
	// A user will grant permission to the app then they will be redirected.
	// That redirect url will contain an authorization code to use to get the user's access token.
	GetAccessTokenFromAuthorizationCode(authorizationCode string) (*TokenResponse, error)
	// Get the access token from strava for an existing application
	//
	// The user must already have granted access to the application to use this.
	// (Otherwise, must authenticate first)
	RefreshAccessToken(refreshToken string) (*TokenResponse, error)
}
// Implementation of the Authenticator interface
type Authorizer struct {
	ClientId string
	ClientSecret string
	RedirectUri string
	Scope string
}
// Format the ApprovalUrlFormat
func (a *Authorizer) GenerateApprovalUrl() string {
	return fmt.Sprintf(approvalUrlFormat, a.ClientId, responseType, a.RedirectUri, approvalPrompt, a.Scope)
}
// This is for the FIRST TIME getting the access token.
//
// A user will grant permission to the app then they will be redirected.
// That redirect url will contain an authorization code to use to get the user's access token.
func (a *Authorizer) GetAccessTokenFromAuthorizationCode(authorizationCode string) (*TokenResponse, error) {
	form := auth.GenerateAuthFormData(a.ClientId, a.ClientSecret, authorizationCode, grantTypeAuthorizationCode)
	// Make the HTTP POST request
	response, err := http.PostForm(authorizationUrl, form)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	// Read response body
	body, err3 := io.ReadAll(response.Body)
	if err3 != nil {
		return nil, err3
	}
	// Load response into struct
	var token TokenResponse
	err4 := json.Unmarshal(body, &token)
	if err4 != nil {
		return nil, err4
	}
	return &token, nil
}
// Get the access token from strava for an existing application
//
// The user must already have granted access to the application to use this.
// (Otherwise, must authenticate first)
func (a *Authorizer) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	// Create the data
	payload, err := auth.MakePayloadRefreshToken(a.ClientId, a.ClientSecret, refreshToken, grantTypeRefreshToken)
	if err != nil {
		return nil, err
	}
	// Create a request
	req, err1 := http.NewRequest("POST", authorizationUrl, bytes.NewBuffer(payload))
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
	var response TokenResponse
	err4 := json.Unmarshal(body, &response)
	if err4 != nil {
		return nil, err4
	}
	return &response, nil
}