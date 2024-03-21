package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	responseType string = "code"
	redirectUri string = "http://localhost/exchange_token"
	approvalPrompt string = "force"
	scope = "activity:read_all"
	ApprovalUrlFormat string = "https://www.strava.com/oauth/authorize?client_id=%s&response_type=%s&redirect_uri=%s&approval_prompt=%s&scope=%s"
)

// Format the ApprovalUrlFormat
func GenerateApprovalUrl() string {
	return fmt.Sprintf(ApprovalUrlFormat, ClientId, responseType, redirectUri, approvalPrompt, scope)
}
// open the approval url in browser
func InitialAuthorizationDirect() {
	approvalUrl := GenerateApprovalUrl()
	openURL(approvalUrl)
}
// This is for the FIRST TIME getting the access token.
//
// A user will grant permission to the app then they will be redirected.
// That redirect url will contain an authorization code to use to get the user's access token.
func GetAccessTokenFromAuthorizationCode(authorizationCode string) (*GetAccessTokenResponse, error) {
	// Prepare form data
	form := url.Values{}
	form.Set("client_id", ClientId)
	form.Set("client_secret", ClientSecret)
	form.Set("code", authorizationCode)
	form.Set("grant_type", "authorization_code")

	// Make the HTTP POST request
	response, err := http.PostForm(AuthorizationUrl, form)
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
	var token GetAccessTokenResponse
	err4 := json.Unmarshal(body, &token)
	if err4 != nil {
		return nil, err4
	}
	return &token, nil
}