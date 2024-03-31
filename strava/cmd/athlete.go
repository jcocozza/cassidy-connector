package cmd

import (
	"github.com/spf13/cobra"
)

var getAthlete = &cobra.Command{
	Use: "athlete",
	Short: "Get an authenticated athlete. Expects an access token.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		accessToken := args[0]
		stravaApp := createApp()
		stravaApp.GetAthlete(accessToken)
	},
}

func init() {
	rootCmd.AddCommand(getAthlete)
}