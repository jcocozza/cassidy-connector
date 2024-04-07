package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.0.1"

var useCassidyApp bool

// app flags
var clientId string
var clientSecret string
var redirectUri string
var scopes []string

var outputPath string

var rootCmd = &cobra.Command{
	Use:   "cassidy-strava",
	Version: version,
	Short: "cassidy-strava is a cli tool to interact with the Strava API",
	Long: `cassidy-strava is a cli tool to interact with the Strava API`,
	Run: func(cmd *cobra.Command, args []string) {
	  // Do Stuff Here
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&useCassidyApp, "use-cassidy", "c", false, "include this flag if you want to use the cassidy's strava application as opposed to your own.")

	rootCmd.PersistentFlags().StringVar(&clientId, "client-id", "", "the client id of your strava application")
	rootCmd.PersistentFlags().StringVar(&clientSecret, "client-secret", "", "the client secret of your strava application")
	rootCmd.PersistentFlags().StringVar(&redirectUri, "redirect-uri", "http://localhost/exchange_token", "the redirect uri of your strava application")
	rootCmd.PersistentFlags().StringSliceVar(&scopes, "scope", []string{"activity:read_all"}, "the scope requirement of your strava application")

	rootCmd.PersistentFlags().StringVarP(&outputPath, "path", "f", "", "The path to save successful output to. (will not write errors at this time)")

	rootCmd.MarkFlagsRequiredTogether("client-id", "client-secret")
	rootCmd.MarkFlagsMutuallyExclusive("use-cassidy", "client-id")
	rootCmd.MarkFlagsMutuallyExclusive("use-cassidy", "client-secret")
	rootCmd.MarkFlagsMutuallyExclusive("use-cassidy", "scope")
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}