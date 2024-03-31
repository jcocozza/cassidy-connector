package app

import (
	"github.com/jcocozza/cassidy-connector/strava/app/api"
	"github.com/jcocozza/cassidy-connector/strava/app/auth"
	config "github.com/jcocozza/cassidy-connector/strava/internal"
)

type StravaApp interface {
	auth.Authenticator
	api.StravaAPI
}
// The implementation of StravaApp interface
type App struct {
	auth.Authorizer
	api.StravaAPICaller
}
func NewApp(clientId string, clientSecret, redirectUri string, scope string) *App {
	return &App{
		Authorizer: auth.Authorizer{
			ClientId: clientId,
			ClientSecret: clientSecret,
			RedirectUri: redirectUri,
			Scope: scope,
		},
		StravaAPICaller: api.StravaAPICaller{},
	}
}
// Create the default Cassidy App for those who don't want to create their own strava app
func CassidyApp(redirectUri string) *App {
	return NewApp(config.ClientId, config.ClientSecret, redirectUri, config.Scope)
}