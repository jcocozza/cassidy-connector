package cmd

import (
	"log/slog"
	"os"

	"github.com/jcocozza/cassidy-connector/strava/app"
	"golang.org/x/oauth2"
)

// Create the app based on the passed flag settings
func createApp() (*app.App, *oauth2.Token, error) {
	var tkn *oauth2.Token
	var err error
	stravaApp := app.NewApp(
		clientId,
		clientSecret,
		authorizationCallbackDomain,
		callbackPath,
		webhookPath,
		webhookServerURL,
		webhookVerifyToken,
		nil,
		scopes,
		//nil,// no logger for the cli
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)
	// when we have a token, we want to load it in to the app
	if tokenPath != "" {
		tkn, err = stravaApp.ReadTokenFromFile(tokenPath)
		if err != nil {
			return nil, nil, err
		}
	} else if token != "" {
		tkn, err = stravaApp.ReadTokenString(token)
		if err != nil {
			return nil, nil, err
		}
	}
	return stravaApp, tkn, nil
}
