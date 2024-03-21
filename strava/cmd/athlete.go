package cmd

import (
	"github.com/jcocozza/cassidy-connector/strava/api"
	"github.com/spf13/cobra"
)

var getAthlete = &cobra.Command{
	Use: "athlete",
	Short: "Get an authenticated athlete. Expects an access token.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		accessToken := args[0]
		api.GetAthlete(accessToken)
	},
}

func init() {
	rootCmd.AddCommand(getAthlete)
}