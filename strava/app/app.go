package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jcocozza/cassidy-connector/strava/app/api"
	"github.com/jcocozza/cassidy-connector/strava/swagger"
	"github.com/jcocozza/cassidy-connector/strava/utils"

	"golang.org/x/oauth2"
)

const (
	responseType            string = "code"
	approvalPrompt          string = "force"
	approvalUrlFormat       string = "https://www.strava.com/oauth/authorize?client_id=%s&response_type=%s&redirect_uri=%s&approval_prompt=%s&scope=%s"
	stravaAppSettings       string = "https://www.strava.com/settings/apps"
	webhookSubscriptionsURL string = "https://www.strava.com/api/v3/push_subscriptions"

	AspectTypeCreate string = "create"
	AspectTypeUpdate string = "update"
	AspectTypeDelete string = "delete"
)

// A StravaEvent is an event that is sent from the webhook
type StravaEvent struct {
	// either "activity" or "athlete"
	ObjectType string `json:"object_type"`
	// activity id or athlete id based on ObjectType
	ObjectID int `json:"object_id"`
	// either "create", "update" or "delete"
	AspectType string `json:"aspect_type"`
	// only for AspectType = "update"
	// possible keys: "title", "type", "private", "authorized"
	Updates map[string]string `json:"updates"`
	// athlete's id
	OwnerID int `json:"owner_id"`
	// push subscription id receiving the event
	SubscriptionID int `json:"subscription_id"`
	// time that the event occured
	EventTime int `json:"event_time"`
}

// An app is a way of interacting with the strava api.
//
// There are two main components to this struct:
//  1. The Strava API application. These are created by strava users and managed at `https://www.strava.com/settings/api`.
//     These are the `ClientId`, `ClientSecret`, `RedirectURL`, `Scopes`.
//     Given these identifiers we can properly interact with the OAuth2 Strava API (which is #2)
//  2. The interaction with the Strava API.
//     This is handled via OAuth2.
//     The `App` struct contains the necessary methods for authenticating and connecting Strava API applications.
//     This is handled by `OAuthConfig`, `SwaggerConfig`, `StravaClient`, and `Token`.
//     `StravaClient` also exposes the various Swagger API services for those that want to use the swagger methods directly.
//     The swagger methods/api calls are wrapped by the custom functions that allow for a layer of abstration to simplify interaction with the strava api.
//     This is all found the the `Api` field of the `App` struct
type App struct {
	logger       *slog.Logger
	ClientId     string
	ClientSecret string
	RedirectURL  string
	// note that this must be a publicly accessible url otherwise Strava cannot make the challenge request to it
	//
	// this must also be set in your strava application
	// see strava.com/settings/api then press edit.
	// this should be the "Authorization Callback Domain"
	AuthorizationCallbackDomain string
	// where you want the webserver to run for the webhooks e.g. http://localhost:8086
	//
	// Traffic from AuthorizationCallbackDomain should be routed to this server
	WebhookServerURL string
	// Token to verify that data coming from the webhook is what you expect it to be
	// Can just be a random string
	WebhookVerifyToken string
	Scopes             []string
	// OAuthConfig handles OAuth and creates the HTTPClient that is used to make requests for the StravaClient
	OAuthConfig *oauth2.Config
	// The SwaggerConfig is passed into the creation of the StravaClient.
	//
	// It is exposed in the app struct because we want to be able to set the HTTPClient that the StravaClient uses
	SwaggerConfig *swagger.Configuration
	// Contains the methods for interacting with the strava API
	StravaClient *swagger.APIClient
	// A way to get the authorization token from the intial authorization process
	// Any calls to the stravaRedirectHandler will push the authorization code to the AuthorizationReciver channel.
	AuthorizationReciever chan string
	// this ensures that the webhook GET request from strava completes before we move forward
	WebhookReciever chan string
	// optional; a user defined function that tells the api how to handle new events
	//
	// *IMPORTANT* this will be called asynchronously with a go func
	// the strava webhook wants a response in less then 2 seconds so all events need to be handled asynchronously
	//
	// A basic WebhookEventHandler might look like:
	//
	//func stravaEventHandler(se app.StravaEvent) {
	//	switch se.AspectType {
	//	case AspectTypeCreate:
	//		fmt.Println("creating")
	//	case AspectTypeUpdate:
	//		fmt.Println("updating")
	//	case AspectTypeDelete:
	//		fmt.Println("deleteing")
	//	}
	//}
	WebhookEventHandler func(StravaEvent)
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

// note that authorizationCallbackDomain, webhookServerURL, and webhookVerifyToken can be empty strings if you aren't interested in webhooks
func NewApp(clientId string, clientSecret, redirectURL string, authorizationCallbackDomain string, webhookServerURL string, webhookVerifyToken string, webhookEventHandler func(StravaEvent), scopes []string, logger *slog.Logger) *App {
	approvalUrl := generateApprovalUrl(clientId, redirectURL, scopes)
	oauthCfg := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  approvalUrl,
			TokenURL: "https://www.strava.com/oauth/token",
		},
	}
	cfg := swagger.NewConfiguration()
	client := swagger.NewAPIClient(cfg)
	reciever := make(chan string)
	webhookReciever := make(chan string, 1)
	if logger == nil {
		logger = NoopLogger()
	}
	logger = logger.WithGroup("cassidy-strava")
	return &App{
		logger:                      logger,
		ClientId:                    clientId,
		ClientSecret:                clientSecret,
		RedirectURL:                 redirectURL,
		AuthorizationCallbackDomain: authorizationCallbackDomain,
		WebhookServerURL:            webhookServerURL,
		WebhookVerifyToken:          webhookVerifyToken,
		WebhookReciever:             webhookReciever,
		WebhookEventHandler:         webhookEventHandler,
		Scopes:                      scopes,
		SwaggerConfig:               cfg,
		OAuthConfig:                 oauthCfg,
		StravaClient:                client,
		Api:                         api.NewStravaAPI(client, logger.WithGroup("api")),
		AuthorizationReciever:       reciever,
	}
}

// Return the approval url
func (a *App) ApprovalUrl() string {
	scopeStr := strings.Join(a.Scopes, ",")
	a.logger.Debug("generating approval url", slog.Any("scope string", scopeStr))
	return fmt.Sprintf(approvalUrlFormat, a.ClientId, responseType, a.RedirectURL, approvalPrompt, scopeStr)
}

// This is for the FIRST TIME getting the access token.
//
// A user will grant permission to the app then will be redirected to the application's RedirectURL.
// The RedirectURL will contain an authorization code. This code is used to get the user's access token.
//
// You are responsible for persisting user tokens
func (a *App) GetAccessTokenFromAuthorizationCode(ctx context.Context, code string) (*oauth2.Token, error) {
	a.logger.InfoContext(ctx, "getting access token from authorization code")
	token, err := a.OAuthConfig.Exchange(ctx, code)
	if err != nil {
		a.logger.ErrorContext(ctx, "token exchange failed", slog.String("error", err.Error()))
		return nil, err
	}
	return token, nil
}

// Turn a json string token into an `oauth2.Token` struct
func (a *App) ReadTokenString(tokenJsonString string) (*oauth2.Token, error) {
	var token *oauth2.Token
	err := json.Unmarshal([]byte(tokenJsonString), &token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Load an oauth2 token from a .json file
func (a *App) ReadTokenFromFile(tokenFilePath string) (*oauth2.Token, error) {
	tokenData, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	err = json.Unmarshal(tokenData, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// Get the authorization code form the url that results from the redirect.
// This is written to the AuthorizationReciever channel.
//
// If there is an error (e.g. the user denies access permission), write the error to the AuthorizationReciver channel with an "error:" prefix.
func (a *App) stravaRedirectHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Debug("strava redirect handler called")
	// Extract URL parameters here and handle them accordingly
	code := r.URL.Query().Get("code") // Assuming 'code' is the parameter sent by Strava
	err := r.URL.Query().Get("error") // if the user denies, the url will send an error "access_denied"
	if code != "" {
		a.logger.Debug("sending code to authorization reciever")
		a.AuthorizationReciever <- code
	} else if err != "" {
		a.logger.Warn("sending error to authorization reciever", slog.String("error", err))
		a.AuthorizationReciever <- "error:" + err
	}
}

// Parse a url into its "address:port" and its "url/path"
//
// e.g. http://localhost:9999/strava/callback -> "localhost:9999", "strava/callback", err
func parseURL(inputURL string) (string, string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", "", err
	}
	if len(parsedURL.Path) > 0 {
		if string(parsedURL.Path[0]) == "/" {
			return parsedURL.Host, parsedURL.Path[1:], nil // [1:] is used to remove the leading '/'
		}
		return parsedURL.Host, parsedURL.Path, nil
	} else {
		return parsedURL.Host, parsedURL.Path, nil
	}
}

// Listen to the redirect route. Once the user is directed to it, we can extract the token from the url.
//
// Returns the Http server instance, so it can be shutdown when you like.
func (a *App) StartStravaHttpServer() (*http.Server, error) {
	hostWithPort, path, err := parseURL(a.RedirectURL)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	srv := &http.Server{Addr: hostWithPort, Handler: mux}
	mux.HandleFunc("/"+path, a.stravaRedirectHandler)
	a.logger.Info("starting strava http server", slog.String("address", hostWithPort))
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			a.logger.Error("strava http server failed", slog.String("error", err.Error()))
		}
	}()
	return srv, nil
}

// Run this function when you send the user to strava's authorization site.
//
// `timeoutDuration` is the time in seconds wait before returning nothing. Use -1 for no timeout duration.
// If the duration is exceeded throws an error.
//
// It will start an http listener that listens on the redirect route provided by your app.
// Once the user authorizes the app and is redirected, the http ListenAndServe will detect the authorization code and push it to the AuthorizationReciever channel.
// Finally, the GetAccessTokenFromAuthorizationCode will set the app's token.
//
// From there, you can persist the token in whatever way you please for further access.
//
// For lower level control over this process, you can start the HttpListener, await the code in the channel and get the access token separately:
/*
	server, err := a.StartStravaHttpServer()
	if err != nil {
		return nil, err
	}
	code := <-s.App.AuthorizationReciever
	fmt.Println("GOT CODE:" + code)
	err := s.App.GetAccessTokenFromAuthorizationCode(context.TODO(), code)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
*/
//
func (a *App) AwaitInitialToken(timeoutDuration int) (*oauth2.Token, error) {
	// TODO: allow the context to be passed in
	ctx := context.TODO()
	a.logger.DebugContext(ctx, "awaiting initial token", slog.Int("timeout duration", timeoutDuration))
	// start listening in a separate go routine
	server, err := a.StartStravaHttpServer()
	if err != nil {
		return nil, err
	}
	if timeoutDuration == -1 {
		code := <-a.AuthorizationReciever
		if strings.Contains(code, "error") {
			server.Shutdown(ctx)
			return nil, fmt.Errorf(code)
		}
		token, err := a.GetAccessTokenFromAuthorizationCode(ctx, code)
		if err != nil {
			a.logger.WarnContext(ctx, "failed to get token, shutting down server")
			server.Shutdown(ctx)
			return nil, err
		}
		a.logger.InfoContext(ctx, "got token, shutting down server")
		server.Shutdown(ctx)
		return token, nil
	} else {
		select {
		case code := <-a.AuthorizationReciever:
			if strings.Contains(code, "error") {
				server.Shutdown(ctx)
				return nil, fmt.Errorf(code)
			}
			// recieved token
			token, err := a.GetAccessTokenFromAuthorizationCode(ctx, code)
			if err != nil {
				a.logger.WarnContext(ctx, "failed to get token, shutting down server")
				server.Shutdown(ctx)
				return nil, err
			}
			server.Shutdown(ctx)
			return token, nil
		case <-time.After(time.Duration(timeoutDuration) * time.Second):
			// didn't recieve token in time
			a.logger.WarnContext(ctx, "hit timeout, failed to get token, shutting down server")
			server.Shutdown(ctx)
			return nil, fmt.Errorf("exceeded timeout duration")
		}
	}
}

// Open the Approval Url in the users browser
func (a *App) OpenAuthorizationGrant() {
	url := a.ApprovalUrl()
	a.logger.Debug("opening authorization grant", slog.String("url", url))
	utils.OpenURL(url)
}

// Open the strava settings page
//
// This idea is to make it easy for the users to deauthenticate/revoke access to the app whenever they like.
func (a *App) OpenStravaAppSettings() {
	a.logger.Debug("opening strava app settings", slog.String("url", stravaAppSettings))
	utils.OpenURL(stravaAppSettings)
}

// this handler does a great deal of work
//
// When a get request is make, (that is creating a subscription to the webhook):
//
//	must respond within 2 seconds to the get request from strava
//	per https://developers.strava.com/docs/webhooks/ it must repond with http status 200 and the hub.challenge
//	once this happens the original webhook POST request will receive a response
//
// When a post request is made strava is sending new events
func (a *App) webhookRedirectHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Debug("webhook redirect handler called")
	switch r.Method {
	case http.MethodGet:
		a.logger.Debug("webhook redirect handler method is GET")
		challenge := r.URL.Query().Get("hub.challenge")
		verificationToken := r.URL.Query().Get("hub.verify_token")
		// if the verification token is not the same as when we created the subscription then we are not receiving the correct response
		if verificationToken != a.WebhookVerifyToken {
			a.logger.Warn("verification tokens do not match")
			http.Error(w, "verification tokens do not match", http.StatusUnauthorized)
			return
		}
		response, err := json.Marshal(map[string]string{"hub.challenge": challenge})
		if err != nil {
			a.logger.Error("hub.challenge could not be marshalled")
			http.Error(w, "hub.challenge could not be marshalled", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		a.logger.Info("sending challenge to webhook reciever")
		a.WebhookReciever <- challenge // send the challenge to let the main request to read the post
	case http.MethodPost:
		a.logger.Debug("webhook redirect handler method is POST")
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			a.logger.Error("unable to read post content")
			http.Error(w, fmt.Sprintf("error reading post content: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		se := StravaEvent{}
		err = json.Unmarshal(body, &se)
		if err != nil {
			a.logger.Error("unable to unmarshall strava event")
			http.Error(w, fmt.Sprintf("error unmarshalling event: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		if a.WebhookEventHandler != nil {
			a.logger.Debug("running webhook event handler")
			go func() {
				a.WebhookEventHandler(se)
			}()
		} else {
			a.logger.Warn("no webhook event handler defined. doing nothing")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		a.logger.Debug("unexpected method. doing nothing.")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// this is a one time run allowing you to subscribe to strava webhooks
// returns:
//   - the subscription id and the server that will be called to get events
//   - the created server
//   - a wait group. by calling wg.Wait() you keep the server running until it is explicitly stopped.
//
// note that the AuthorizationCallbackDomain MUST be open to the internet otherwise strava cannot send information to the server
func (a *App) CreateSubscription() (int, *http.Server, *sync.WaitGroup, error) {
	a.logger.Debug("creating subscription")
	srv, wg, err := a.LaunchWebhookServer()
	if err != nil {
		wg.Done()
		a.logger.Error("launching server failed, unable to create subscription")
		return -1, nil, nil, err
	}
	// the subscription process will return a challence to the callback url
	// as such we need to be listening for that before we make the request
	// make sure that the server is running
	time.Sleep(1 * time.Second)
	payload := map[string]string{
		"client_id":     a.ClientId,
		"client_secret": a.ClientSecret,
		"callback_url":  a.AuthorizationCallbackDomain,
		"verify_token":  a.WebhookVerifyToken,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		wg.Done()
		a.logger.Error("unable to marshal payload. unable to create subscription")
		return -1, nil, nil, err
	}
	req, err := http.NewRequest(http.MethodPost, webhookSubscriptionsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		wg.Done()
		a.logger.Error("unable to create request. unable to create subscription")
		return -1, nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		wg.Done()
		a.logger.Error("request failed to run. unable to create subscription")
		return -1, nil, nil, err
	}
	defer resp.Body.Close()
	// wait for the webhook verification challege to complete
	// once this happens, strava has confirmed the webhook, so we are now expecting the response
	a.logger.Debug("awaiting challenge")
	_ = <-a.WebhookReciever
	a.logger.Debug("got challenge")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		wg.Done()
		a.logger.Error("unable to read response body. unable to create subscription")
		return -1, nil, nil, err
	}
	type subscriptionResponse struct {
		Id int `json:"id"`
	}
	sr := subscriptionResponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		wg.Done()
		a.logger.Error("unable to unmarshal subscription response. unsure if subscription has actually been created. investigate further.")
		return -1, nil, nil, err
	}
	a.logger.Info("subscription created", slog.Int("subscription id", sr.Id))
	return sr.Id, srv, wg, nil
}

// spawns a server with 2 routes:
//  1. the path specified by the AuthorizationCallbackDomain which will process events
//  2. /status which will return "alive" if the server is alive
//
// returns:
//   - the created server
//   - a wait group. by calling wg.Wait() you keep the server running until it is explicitly stopped.
func (a *App) LaunchWebhookServer() (*http.Server, *sync.WaitGroup, error) {
	_, path, err := parseURL(a.AuthorizationCallbackDomain)
	if err != nil {
		return nil, nil, err
	}
	hostWithPort, _, err := parseURL(a.WebhookServerURL)
	if err != nil {
		return nil, nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/"+path, a.webhookRedirectHandler)
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("alive")) })
	srv := &http.Server{Addr: hostWithPort, Handler: mux}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	a.logger.Info("launching webhook server", slog.String("address", hostWithPort), slog.String("webhook path", path))
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			a.logger.Error("webhook server failed", slog.String("error", err.Error()))
		}
	}()
	return srv, wg, nil
}

// view the subscription associated with your client id/client secret
// right now, just prints the response
func (a *App) ViewSubscription() error {
	url := fmt.Sprintf(webhookSubscriptionsURL+"?client_id=%s&client_secret=%s", a.ClientId, a.ClientSecret)
	a.logger.Debug("viewing subscription", slog.String("url", url))
	resp, err := http.Get(url)
	if err != nil {
		a.logger.Error("error making request", slog.String("error", err.Error()))
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Error("unable to read response", slog.String("error", err.Error()))
		return err
	}
	// TODO: struct-ify the response
	fmt.Println(string(body))
	return nil
}

// delete the subscription associated with your client id/client secret
func (a *App) DeleteSubscription(subscriptionID string) error {
	url := fmt.Sprintf(webhookSubscriptionsURL+"/%s?client_id=%s&client_secret=%s", subscriptionID, a.ClientId, a.ClientSecret)
	a.logger.Debug("delete subscription", slog.String("subscription id", subscriptionID), slog.String("url", url))
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		a.logger.Error("error creating request", slog.String("error", err.Error()))
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		a.logger.Error("error making request", slog.String("error", err.Error()))
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Error("unable to read response", slog.String("error", err.Error()))
		return err
	}
	// TODO: struct-ify the response
	fmt.Println(string(body))
	return nil
}
