package config
// The following will be unnecessary for most use cases since this tool provides a way for users to use their own strava applications.
// However, if you feel like you want to create an entirely distinct tool for your own purposes then read on...
// For those wishing to distribute their own versions of this package, but wanting to keep their strava application information private should do the following:
// 	1. rename this file from `config.tmpl.go` to `config.go`.
//	2. rename each of the const and variables by removing the `_tmpl` from them. e.g. `ClientId_tmpl` -> ClientId
//	3. Include your own information.
// 	4. Package the config file with a binary that you distribute. This way the application information remains private, but users can access your strava application.
//	5. You may want to rename functions like "CassidyApp()" or variables like "useCassidyApp" as this will be unclear.
//		If you have specified your own variables, then the Cassidy strava application will no longer be used.

const (
	ClientId_tmpl     string = ""
	ClientSecret_tmpl string = ""
	RefreshToken_tmpl string = ""
	RedirectURL_tmpl  string = "http://localhost/exchange_token"
)
var (
	Scopes_tmpl = []string{"activity:read_all"}
)