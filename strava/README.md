# Strava

The strava package provides a CLI tool for interacting with the strava API.
Further, the package exposes methods that allow you to develop your own ways to interact with the strava API.
If you like, you can also work directly with the swagger api call methods for less abstraction.

For each of these interaction methods you can choose to use to work with the strava api through your own strava api app, or the one cassidy provides.

## CLI
The CLI is a relatively easy way to get off the ground and test various things.
```
cassidy-strava is a cli tool to interact with the Strava API

Usage:
  cassidy-strava [flags]
  cassidy-strava [command]

Available Commands:
  api            all subcommands here require a token for authentication
  approval-url   Generate the approval url for the user to grant access.
  completion     Generate the autocompletion script for the specified shell
  help           Help about any command
  initial-access For getting the user's access token for the first time.
  open-grant     Open a browser to grant allow for permission granting

Flags:
      --client-id string       the client id of your strava application
      --client-secret string   the client secret of your strava application
      --config string          the config file of the application. see config.tmpl.json for format. a config is NOT required if you want to pass everything manually. (default is $HOME/.cassidy-connector-strava.json)
  -h, --help                   help for cassidy-strava
  -f, --path string            the path to save successful output to. (will not write errors at this time)
      --redirect-url string    the redirect url of your strava application (default "http://localhost/exchange_token")
      --scopes strings         the scope requirement of your strava application (default [activity:read_all])
  -v, --version                version for cassidy-strava

Use "cassidy-strava [command] --help" for more information about a command.
```

By default, the CLI will check for a config in `$HOME/.cassidy-connector-strava.json`. You can override this location with the `--config` flag.
The config is a JSON file structured as follows:
```
{
    "client_id": "",
    "client_secret": "",
    "redirect_url": "",
    "scopes": [],
    "token_path": ""
}
```
As you can see, there are corresponding flags in the CLI for each of these. The config merely allows the flags to be set when the CLI is run without having to manually enter them each time.

### CLI API
The CLI API command exposes all the commands relevant for interacting with strava data.
```
all subcommands here require a token for authentication

Usage:
  cassidy-strava api [flags]
  cassidy-strava api [command]

Available Commands:
  activities  Get activities.
  activity    Get an activity by activity id. Expects an activity id.
  athlete     Get an authenticated athlete.
  streams     Get streams for a given activity

Flags:
  -h, --help                      help for api
      --token                     a json token. you must include the entire token wrapped in . the json token conforms to `oauth2.Token` struct found here: https://pkg.go.dev/golang.org/x/oauth2#Token. (this is an ugly, but can be useful for testing purposes)
      --token-path oauth2.Token   the path to a .json file that contains an OAuth2 token. This json must conform to the oauth2.Token struct found here: https://pkg.go.dev/golang.org/x/oauth2#Token.

Global Flags:
      --client-id string       the client id of your strava application
      --client-secret string   the client secret of your strava application
      --config string          the config file of the application. see config.tmpl.json for format. a config is NOT required if you want to pass everything manually. (default is $HOME/.cassidy-connector-strava.json)
  -f, --path string            the path to save successful output to. (will not write errors at this time)
      --redirect-url string    the redirect url of your strava application (default "http://localhost/exchange_token")
      --scopes strings         the scope requirement of your strava application (default [activity:read_all])

Use "cassidy-strava api [command] --help" for more information about a command.
```

Importantly, you must obtain an OAuth2 token to have access to the data. As previously mentioned, this token can either be passed into the CLI directly, or can be stored in a file and loaded through the `--token-path` flag, or via config file specification.

## IMPORTANT NOTICE
You may need to change the `LatLng` struct in the `strava/internal/swagger/model_lat_lng.go` file to be a list of `float32` (or `float64`). It appears that the `strava/internal/swagger/make.sh` using `swagger-codegen` generates this improperly.
`type LatLng struct {}` is **INCORRECT**.

The file should look like THIS:
```
/*
 * Strava API v3
 *
 * The [Swagger Playground](https://developers.strava.com/playground) is the easiest way to familiarize yourself with the Strava API by submitting HTTP requests and observing the responses before you write any client code. It will show what a response will look like with different endpoints depending on the authorization scope you receive from your athletes. To use the Playground, go to https://www.strava.com/settings/api and change your “Authorization Callback Domain” to developers.strava.com. Please note, we only support Swagger 2.0. There is a known issue where you can only select one scope at a time. For more information, please check the section “client code” at https://developers.strava.com/docs.
 *
 * API version: 3.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

// A pair of latitude/longitude coordinates, represented as an array of 2 floating point numbers.
type LatLng []float32

```