package auth

import (
	"encoding/json"
	"net/url"
)

type GrantType string

// Create auth form data for authorization code
func GenerateAuthFormData(clientId string, clientSecret string, authorizationCode string, grantType GrantType) url.Values {
	form := url.Values{}
	form.Set("client_id", clientId)
	form.Set("client_secret", clientSecret)
	form.Set("code", authorizationCode)
	form.Set("grant_type", string(grantType))
	return form
}
// Create the payload for a refresh token
func MakePayloadRefreshToken(clientId string, clientSecret string, refreshToken string, grantType GrantType) ([]byte, error) {
	payload := map[string]string{
		"client_id":     clientId,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
		"grant_type":    string(grantType),
		"f":             "json",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return payloadBytes, nil
}