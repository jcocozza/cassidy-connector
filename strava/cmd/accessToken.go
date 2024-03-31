package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)
// Used to get the first access token.
// Users grant permission and strava will return an authorization code.
// This authorization code is used to get the access token and refresh token
var initialAccess = &cobra.Command{
	Use: "initial-access",
	Short: "For getting the user's access token for the first time. Expects an authorization code",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		authorizationCode := args[0]
		stravaApp := createApp()

		token, err := stravaApp.GetAccessTokenFromAuthorizationCode(authorizationCode)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("access token: " + token.AccessToken)
	},
}
// Used to refresh the access token for existing users who have granted access.
var refreshAccessToken = &cobra.Command{
    Use: "refresh",
    Short: "Refresh an access token. Expects a refresh token.",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
		refreshToken := args[0]
		stravaApp := createApp()

		token, err := stravaApp.RefreshAccessToken(refreshToken)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("access token: " + token.AccessToken)
    },
}
func init() {
	rootCmd.AddCommand(initialAccess)
    rootCmd.AddCommand(refreshAccessToken)
}