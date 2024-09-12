package cmd

import "github.com/jcocozza/cassidy-connector/finalSurge/app"

func createApp(email, password string) *app.App {
	return app.NewApp(email, password)
}
