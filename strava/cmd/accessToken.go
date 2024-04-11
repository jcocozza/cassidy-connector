package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// Used to get the first access token.
// Users grant permission and strava will return an authorization code.
// This authorization code is used to get the access token and refresh token
// Your application is responsible for persisting the returned token.
var initialAccess = &cobra.Command{
	Use: "initial-access [authorization code]",
	Short: "For getting the user's access token for the first time.",
	Long: "Used for getting the user's access token for the first time. You are responsible for persisting the returned token so that it can be used later.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		authorizationCode := args[0]
		stravaApp, err := createApp()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err1 := stravaApp.GetAccessTokenFromAuthorizationCode(context.TODO(), authorizationCode)
		if err1 != nil {
			fmt.Println(err1.Error())
			return
		}

		jsonBytes, err := json.Marshal(stravaApp.Token)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if outputPath != "" {
			writeOutput(outputPath, jsonBytes)
		}
		fmt.Println(string(jsonBytes))
	},
}
// Used to get the approval url for the user to grant access
var approvalUrl = &cobra.Command{
	Use: "approval-url",
	Short: "Generate the approval url for the user to grant access.",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		stravaApp, err := createApp()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		url := stravaApp.ApprovalUrl()

		if outputPath != "" {
			writeOutput(outputPath, []byte(url))
		}
		fmt.Println(url)
	},
}

func init() {
	rootCmd.AddCommand(initialAccess)
	rootCmd.AddCommand(approvalUrl)
}