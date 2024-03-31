package cmd

import "github.com/jcocozza/cassidy-connector/strava/app"

// Create the app based on the passed flag settings
func createApp() app.StravaApp {
	var stravaApp *app.App
	if useCassidyApp {
		stravaApp = app.CassidyApp(redirectUri)
	} else {
		stravaApp = app.NewApp(clientId, clientSecret, redirectUri, scope)
	}
	return stravaApp
}