package cmd

import (
	"os"

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
// write a slice of bytes to a file
func writeOutput(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0644)
}