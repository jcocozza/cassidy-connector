package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jcocozza/cassidy-connector/strava/app/api"
	config "github.com/jcocozza/cassidy-connector/strava/internal"
	"github.com/jcocozza/cassidy-connector/strava/internal/swagger"

	"golang.org/x/oauth2"
)

const (
	responseType string = "code"
	approvalPrompt string = "force"
	approvalUrlFormat string = "https://www.strava.com/oauth/authorize?client_id=%s&response_type=%s&redirect_uri=%s&approval_prompt=%s&scope=%s"
)
// An app is a way of interacting with the strava api.
//
// There are two main components to this struct:
// 	1) The Strava API application. These are created by strava users and managed at `https://www.strava.com/settings/api`.
//		These are the `ClientId`, `ClientSecret`, `RedirectURL`, `Scopes`.
//		Given these identifiers we can properly interact with the OAuth2 Strava API (which is #2)
//	2) The interaction with the Strava API.
//		This is handled via OAuth2.
//		The `App` struct contains the necessary methods for authenticating and connecting Strava API applications.
//		This is handled by `OAuthConfig`, `SwaggerConfig`, `StravaClient`, and `Token`.
//		`StravaClient` also exposes the various Swagger API services for those that want to use the swagger methods directly.
//		The swagger methods/api calls are wrapped by the custom functions that allow for a layer of abstration to simplify interaction with the strava api.
//		This is all found the the `Api` field of the `App` struct
type App struct {
	ClientId string
	ClientSecret string
	RedirectURL string
	Scopes []string
	// OAuthConfig handles OAuth and creates the HTTPClient that is used to make requests for the StravaClient
	OAuthConfig *oauth2.Config
	// The SwaggerConfig is passed into the creation of the StravaClient.
	//
	// It is exposed in the app struct because we want to be able to set the HTTPClient that the StravaClient uses
	SwaggerConfig *swagger.Configuration
	// Contains the methods for interacting with the strava API
	StravaClient *swagger.APIClient
	// This contains the user's short-lived access token which is used to access data.
	// When it expires, use the user's refresh token to get a new access token.
	//
	// This struct is obtained in 1 of 2 ways:
	//   - by possessing an existing refresh token and getting a new access token (handled automatically by the oauth2 package).
	//   - via user authorization, whereby an auth code is issued and is used to get the access token.
	Token *oauth2.Token
	// This is where the data methods are called from.
	// It is a layer of abstraction to simplify making calls to the strava API.
	// This is the primary purpose of this package.
	Api *api.StravaAPI
}
// Format the ApprovalUrlFormat
func generateApprovalUrl(clientId string, redirectUrl string, scopes []string) string {
	scopeStr := strings.Join(scopes, ",")
	return fmt.Sprintf(approvalUrlFormat, clientId, responseType, redirectUrl, approvalPrompt, scopeStr)
}
func NewApp(clientId string, clientSecret, redirectURL string, scopes []string) *App {
	approvalUrl := generateApprovalUrl(clientId, redirectURL, scopes)
	oauthCfg := &oauth2.Config{
		ClientID: clientId,
		ClientSecret: clientSecret,
		RedirectURL: redirectURL,
		Scopes: scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL: approvalUrl,
			TokenURL: "https://www.strava.com/oauth/token",
		},
	}
	cfg := swagger.NewConfiguration()
	client := swagger.NewAPIClient(cfg)
	return &App{
		ClientId: clientId,
		ClientSecret: clientSecret,
		RedirectURL: redirectURL,
		Scopes: scopes,

		SwaggerConfig: cfg,
		OAuthConfig: oauthCfg,
		StravaClient: client,
		Api: api.NewStravaAPI(client),
	}
}
// Create the default Cassidy App for those who don't want to create their own strava app
func CassidyApp(redirectURL string) *App {
	return NewApp(config.ClientId, config.ClientSecret, redirectURL, config.Scopes)
}
// Return the approval url
func (a *App) ApprovalUrl() string {
	scopeStr := strings.Join(a.Scopes, ",")
	return fmt.Sprintf(approvalUrlFormat, a.ClientId, responseType, a.RedirectURL, approvalPrompt, scopeStr)
}
// This is for the FIRST TIME getting the access token.
//
// A user will grant permission to the app then will be redirected to the application's RedirectURL.
// The RedirectURL will contain an authorization code. This code is used to get the user's access token.
func (a *App) GetAccessTokenFromAuthorizationCode(ctx context.Context, code string) error {
	token, err := a.OAuthConfig.Exchange(ctx, code)
    if err != nil {
        return err
    }
	httpClient := a.OAuthConfig.Client(ctx, token)
	a.Token = token
	a.SwaggerConfig.HTTPClient = httpClient
	return nil
}
// Turn a json string token into an `oauth2.Token` struct and load it into the app
func (a *App) LoadToken(tokenJsonString string) error {
	var token oauth2.Token
	err := json.Unmarshal([]byte(tokenJsonString), &token)
    if err != nil {
        return err
    }
	httpClient := a.OAuthConfig.Client(context.TODO(), &token)
	a.Token = &token
	a.SwaggerConfig.HTTPClient = httpClient
	return nil
}