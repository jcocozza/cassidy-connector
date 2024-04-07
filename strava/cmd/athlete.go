package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var getAthlete = &cobra.Command{
	Use: "athlete [access token]",
	Short: "Get an authenticated athlete.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		accessTokenString := args[0]
		stravaApp := createApp()
		err := stravaApp.LoadToken(accessTokenString)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		athlete, err := stravaApp.Api.GetAthlete(context.TODO())
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		athleteJsonBytes, err := json.Marshal(athlete)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if outputPath != "" {
			writeOutput(outputPath, athleteJsonBytes)
		}
		fmt.Println(string(athleteJsonBytes))
	},
}

func init() {
	rootCmd.AddCommand(getAthlete)
}