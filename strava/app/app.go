package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jcocozza/cassidy-connector/strava/app/api"
	//config "github.com/jcocozza/cassidy-connector/strava/internal"
	"github.com/jcocozza/cassidy-connector/strava/internal/swagger"
	"github.com/jcocozza/cassidy-connector/strava/utils"

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
	// A way to get the authorization token from the intial authorization process
	AuthorizationReciever chan string
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
	reciever := make(chan string)
	return &App{
		ClientId: clientId,
		ClientSecret: clientSecret,
		RedirectURL: redirectURL,
		Scopes: scopes,

		SwaggerConfig: cfg,
		OAuthConfig: oauthCfg,
		StravaClient: client,
		Api: api.NewStravaAPI(client),
		AuthorizationReciever: reciever,
	}
}
/*
// Create the default Cassidy App for those who don't want to create their own strava app
func CassidyApp(redirectURL string) *App {
	return NewApp(config.ClientId, config.ClientSecret, redirectURL, config.Scopes)
}
*/
// Return the approval url
func (a *App) ApprovalUrl() string {
	scopeStr := strings.Join(a.Scopes, ",")
	return fmt.Sprintf(approvalUrlFormat, a.ClientId, responseType, a.RedirectURL, approvalPrompt, scopeStr)
}
// This is for the FIRST TIME getting the access token. It will set the token internally to the app.
//
// A user will grant permission to the app then will be redirected to the application's RedirectURL.
// The RedirectURL will contain an authorization code. This code is used to get the user's access token.
func (a *App) GetAccessTokenFromAuthorizationCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := a.OAuthConfig.Exchange(ctx, code)
    if err != nil {
        return nil, err
    }
	httpClient := a.OAuthConfig.Client(ctx, token)
	a.Token = token
	a.SwaggerConfig.HTTPClient = httpClient
	return token, nil
}
// Turn a json string token into an `oauth2.Token` struct and load it into the app
func (a *App) LoadTokenString(tokenJsonString string) error {
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
// Load an oauth2 token into the app
func (a *App) LoadTokenDirect(token *oauth2.Token) {
	httpClient := a.OAuthConfig.Client(context.TODO(), token)
	a.Token = token
	a.SwaggerConfig.HTTPClient = httpClient
}
// Load an oauth2 token into the app from a .json file
func (a *App) LoadTokenFromFile(tokenFilePath string) error {
	tokenData, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return err
	}

	var token oauth2.Token
	err = json.Unmarshal(tokenData, &token)
	if err != nil {
		return err
	}

	a.LoadTokenDirect(&token)
	return nil
}
// Get the authorization code form the url that results from the redirect
// TODO: Handle denial of permission / error
func (a *App) stravaRedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract URL parameters here and handle them accordingly
	code := r.URL.Query().Get("code") // Assuming 'code' is the parameter sent by Strava
	fmt.Println(code)
	a.AuthorizationReciever <- code
}
// Parse a url into its "address:port" and its "url/path"
//
// e.g. http://localhost:9999/strava/callback -> "localhost:9999", "strava/callback", err
func parseURL(inputURL string) (string, string, error) {
    parsedURL, err := url.Parse(inputURL)
    if err != nil {
        return "", "", err
    }
    return parsedURL.Host, parsedURL.Path[1:], nil // [1:] is used to remove the leading '/'
}
// Listen to the redirect route. Once the user is directed to it, we can extract the token from the url.
func (a *App) StartStravaHttpListener() error {
	hostWithPort, path, err := parseURL(a.RedirectURL)
	if err != nil {
		return err
	}

	fmt.Println("Running ListenAndServe on: " + hostWithPort + " at path: /" + path)
	http.HandleFunc("/" + path, a.stravaRedirectHandler)
	return http.ListenAndServe(hostWithPort, nil)
}
// Run this function when you send the user to strava's authorization site.
//
// It will start an http listener that listens on the redirect route provided by your app.
// Once the user authorizes the app and is redirected, the http ListenAndServe will detect the authorization code and push it to the AuthorizationReciever channel.
// Finally, the GetAccessTokenFromAuthorizationCode will set the app's token.
//
// From there, you can persist the token in whatever way you please for further access.
//
// TODO: Add a timeout to this
func (a *App) AwaitInitialToken() *oauth2.Token {
	go func() {
		err := a.StartStravaHttpListener()
		if err != nil {
			fmt.Println("ListenAndServe: ", err.Error())
		}
	}()

	code := <- a.AuthorizationReciever
	token, err := a.GetAccessTokenFromAuthorizationCode(context.TODO(), code)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return token
}
// Open the Approval Url in the users browser
func (a *App) OpenAuthorizationGrant() {
	url := a.ApprovalUrl()
	utils.OpenURL(url)
}
// Create the OAuth2 token that is used for authentication in the app.
//
// The primary usecase for this is reading in a saved token from a database or file.
// Once you've read in the token information, you can easily create a token with this method.
// Then you can load the token into the app via the `LoadTokenDirect()` function.
func (a *App) createToken(accessToken string, tokenType string, refreshToken string, expiry time.Time) *oauth2.Token {
	return &oauth2.Token{
		AccessToken: accessToken,
		TokenType: tokenType,
		RefreshToken: refreshToken,
		Expiry: expiry,
	}
}