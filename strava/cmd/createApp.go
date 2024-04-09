package cmd

import (
	"os"

	"github.com/jcocozza/cassidy-connector/strava/app"
)

// Create the app based on the passed flag settings
func createApp() *app.App {
	var stravaApp *app.App
	if useCassidyApp {
		stravaApp = app.CassidyApp(redirectURL)
	} else {
		stravaApp = app.NewApp(clientId, clientSecret, redirectURL, scopes)
	}
	return stravaApp
}
// write a slice of bytes to a file
func writeOutput(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0644)
}