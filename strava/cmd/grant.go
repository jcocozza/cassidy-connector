package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jcocozza/cassidy-connector/strava/utils"
)

// Will open the browser for the stored application in the config file.
var grantPermission = &cobra.Command{
	Use: "open-grant",
	Short: "Open a browser to grant allow for permission granting",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		stravaApp := createApp()
		url := stravaApp.ApprovalUrl()
		utils.OpenURL(url)
	},
}

func init() {
	rootCmd.AddCommand(grantPermission)
}
