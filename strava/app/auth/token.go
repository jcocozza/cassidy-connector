package auth

import (
	"github.com/jcocozza/cassidy-connector/strava/internal/auth"
)

const (
	authorizationUrl       string = "https://www.strava.com/oauth/token" // strava's auth endpoint
	grantTypeRefreshToken  auth.GrantType = "refresh_token"
	grantTypeAuthorizationCode auth.GrantType = "authorization_code"
)
// This contains the user's short-lived access token which is used to access data.
// When it expires, use the user's refresh token to get a new access token.
//
// This struct is obtained in 1 of 2 ways:
//   - by possessing an existing refresh token and getting a new access token.
//   - via user authorization, whereby an auth code is issued which is used to get the access token.
type TokenResponse struct {
	TokenType    string `json:"token_type"` // the type of token returned
	AccessToken  string `json:"access_token"` // the access token
	ExpiresAt    int    `json:"expires_at"` // the time that the access token expires
	ExpiresIn    int    `json:"expires_in"` // how long until the access token expires
	RefreshToken string `json:"refresh_token"` // the user's refresh token
}