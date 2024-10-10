package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	version       string = "0.0.1"
	defaultConfig string = ".cassidy-connector-strava.json"
)

type cfg struct {
	ClientId                    string   `json:"client_id"`
	ClientSecret                string   `json:"client_secret"`
	AuthorizationCallbackDomain string   `json:"authorization_callback_domain"`
	CallbackPath				string 	 `json:"callback_path"`
	WebhookPath                 string 	 `json:"webhook_path"`
	WebhookServerURL            string   `json:"webhook_server_url"`
	WebhookVerifyToken          string   `json:"webhook_verify_token"`
	Scopes                      []string `json:"scopes"`
	TokenPath                   string   `json:"token_path"`
}

// global app flag variables
var configPath string
var tokenPath string
var token string
var clientId string
var clientSecret string
var redirectURL string
var authorizationCallbackDomain string
var callbackPath string
var webhookPath string
var webhookServerURL string
var webhookVerifyToken string
var scopes []string
var outputPath string

var RootCmd = &cobra.Command{
	Use:     "cassidy-strava",
	Version: version,
	Short:   "cassidy-strava is a cli tool to interact with the Strava API",
	Long:    `cassidy-strava is a cli tool to interact with the Strava API`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var tokenCmdGroup = &cobra.Command{
	Use:   "api",
	Short: "all subcommands here require a token for authentication",
	Run: func(cmd *cobra.Command, args []string) {
		// Nothing to see here
	},
}

func initConfig() {
	finalConfigPath := ""
	if configPath != "" {
		finalConfigPath = configPath
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		defaultConfigPath := home + "/" + defaultConfig
		// Check if the file exists
		if _, err := os.Stat(defaultConfigPath); err == nil {
			// config exists
			finalConfigPath = defaultConfigPath
		} else if os.IsNotExist(err) {
			// config file does not exist - so there is nothing to do, simply return
			return
		} else {
			// unknown error
			cobra.CheckErr(fmt.Errorf("uh oh, an unknown error occured"))
		}
	}
	data, err := os.ReadFile(finalConfigPath)
	cobra.CheckErr(err)
	var config cfg
	err = json.Unmarshal(data, &config)
	cobra.CheckErr(err)
	// set the relevant flags based on what the config provides
	if config.ClientId != "" {
		RootCmd.Flags().Set("client-id", config.ClientId)
	}
	if config.ClientSecret != "" {
		RootCmd.Flags().Set("client-secret", config.ClientSecret)
	}
	if config.AuthorizationCallbackDomain != "" {
		RootCmd.Flags().Set("auth-callback-domain", config.AuthorizationCallbackDomain)
	}
	if config.CallbackPath != "" {
		RootCmd.Flags().Set("callback-path", config.CallbackPath)
	}
	if config.WebhookPath != "" {
		RootCmd.Flags().Set("webhook-path", config.WebhookPath)
	}
	if config.WebhookServerURL != "" {
		RootCmd.Flags().Set("webhook-server-url", config.WebhookServerURL)
	}
	if config.WebhookVerifyToken != "" {
		RootCmd.Flags().Set("webhook-verify-token", config.WebhookVerifyToken)
	}
	if len(config.Scopes) > 0 {
		scopeStr := strings.Join(config.Scopes, ",")
		RootCmd.Flags().Set("scopes", scopeStr)
	}
	if config.TokenPath != "" {
		tokenCmdGroup.Flags().Set("token-path", config.TokenPath)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&configPath, "config", "", fmt.Sprintf("the config file of the application. see config.tmpl.json for format. a config is NOT required if you want to pass everything manually. (default is $HOME/%s)", defaultConfig))
	RootCmd.PersistentFlags().StringVar(&clientId, "client-id", "", "the client id of your strava application")
	RootCmd.PersistentFlags().StringVar(&clientSecret, "client-secret", "", "the client secret of your strava application")
	RootCmd.PersistentFlags().StringVar(&authorizationCallbackDomain, "auth-callback-domain", "", "the redirect url for the webhook. MUST be exposed to the internet and all traffic from here must be routed to webhook-server-url")
	RootCmd.PersistentFlags().StringVar(&callbackPath, "callback-path", "/strava/callback", "the path off the authorization callback domain for the auth callback")
	RootCmd.PersistentFlags().StringVar(&webhookPath, "webhook-path", "/webhook", "the path off the authorization callback domain for the webhook")
	RootCmd.PersistentFlags().StringVar(&webhookServerURL, "webhook-server-url", "http://localhost:8086", "where the webhook server will run from")
	RootCmd.PersistentFlags().StringVar(&webhookVerifyToken, "webhook-verify-token", "STRAVA", "the verification token for the webook just an arbritary token for verifying your app is receiving the correct webhook responses")
	RootCmd.PersistentFlags().StringSliceVar(&scopes, "scopes", []string{"activity:read_all"}, "the scope requirement of your strava application")
	RootCmd.PersistentFlags().StringVarP(&outputPath, "path", "f", "", "the path to save successful output to. (will not write errors at this time)")
	RootCmd.MarkFlagsRequiredTogether("client-id", "client-secret")
	tokenCmdGroup.PersistentFlags().StringVar(&tokenPath, "token-path", "", "the path to a .json file that contains an OAuth2 token. This json must conform to the `oauth2.Token` struct found here: https://pkg.go.dev/golang.org/x/oauth2#Token.")
	tokenCmdGroup.PersistentFlags().StringVar(&token, "token", "", "a json token. you must include the entire token wrapped in ``. the json token conforms to `oauth2.Token` struct found here: https://pkg.go.dev/golang.org/x/oauth2#Token. (this is an ugly, but can be useful for testing purposes)")
	tokenCmdGroup.MarkFlagsMutuallyExclusive("token-path", "token")
	tokenCmdGroup.MarkFlagsOneRequired("token", "token-path")
	RootCmd.AddCommand(tokenCmdGroup)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
