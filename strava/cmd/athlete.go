package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jcocozza/cassidy-connector/utils"
	"github.com/spf13/cobra"
)

var getAthlete = &cobra.Command{
	Use: "athlete",
	Short: "Get an authenticated athlete.",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		stravaApp, tkn, err := createApp()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		athlete, err := stravaApp.Api.GetAthlete(context.TODO(), tkn)
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
			utils.WriteOutput(outputPath, athleteJsonBytes)
		}
		fmt.Println(string(athleteJsonBytes))
	},
}

func init() {
	tokenCmdGroup.AddCommand(getAthlete)
}
