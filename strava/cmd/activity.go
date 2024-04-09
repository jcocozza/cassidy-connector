package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var includeAllEfforts bool
var getActivity = &cobra.Command{
	Use: "activity [access token] [activity id]",
	Short: "Get an activity by activity id. Expects access token and activity id.",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tokenString := args[0]
		stravaApp := createApp()
		err := stravaApp.LoadToken(tokenString)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		idString := args[1]
		activityId, err := strconv.Atoi(idString)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		activity, err := stravaApp.Api.GetActivity(context.TODO(), activityId, includeAllEfforts)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		activityJsonBytes, err := json.Marshal(activity)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if outputPath != "" {
			writeOutput(outputPath, activityJsonBytes)
		}
		fmt.Println(string(activityJsonBytes))
	},
}
func init() {
	getActivity.Flags().BoolVarP(&includeAllEfforts, "include-all-efforts", "i", false, "include all segment efforts")

	rootCmd.AddCommand(getActivity)
}