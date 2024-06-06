package cmd

import (
	"github.com/jcocozza/cassidy-connector/strava/app"
)

// Create the app based on the passed flag settings
func createApp() (*app.App, error) {
	stravaApp := app.NewApp(clientId, clientSecret, redirectURL, scopes)
	// when we have a token, we want to load it in to the app
	if tokenPath != "" {
		err := stravaApp.LoadTokenFromFile(tokenPath)
		if err != nil {
			return nil, err
		}
	} else if token != "" {
		err := stravaApp.LoadTokenString(token)
		if err != nil {
			return nil, err
		}
	}

	return stravaApp, nil
}
