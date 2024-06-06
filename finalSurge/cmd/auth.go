package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jcocozza/cassidy-connector/utils"

	"github.com/spf13/cobra"
)

var authenticate = &cobra.Command{
	Use:   "authenticate",
	Short: "Authenticate the app",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		finalSurgeApp := createApp(email, password)
		auth, err := finalSurgeApp.Authenticate(context.TODO())
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		authBytes, err := json.Marshal(auth)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if outputPath != "" {
		 	utils.WriteOutput(outputPath, authBytes)
		}
		fmt.Println(string(authBytes))
	},
}

func init() {
	RootCmd.AddCommand(authenticate)
}