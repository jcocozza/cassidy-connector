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
	Use: "activity [activity id]",
	Short: "Get an activity by activity id. Expects an activity id.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stravaApp, err := createApp()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		idString := args[0]
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
	tokenCmdGroup.AddCommand(getActivity)
}