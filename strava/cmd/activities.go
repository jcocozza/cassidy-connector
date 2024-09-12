package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jcocozza/cassidy-connector/utils"
	"github.com/spf13/cobra"
)

const layout string = "2006-01-02"
const layoutInterpretation string = "YYYY-MM-DD"
var perPage int
var before string
var after string
var getActivities = &cobra.Command{
	Use: "activities",
	Short: "Get activities.",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		stravaApp, err := createApp()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var beforeTimePtr *time.Time = nil
		var afterTimePtr *time.Time = nil
		if before != "" {
			beforeTime, err := time.Parse(layout, before)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			beforeTimePtr = &beforeTime
		}
		if after != "" {
			afterTime, err := time.Parse(layout, after)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			afterTimePtr = &afterTime
		}

		activities, err := stravaApp.Api.GetActivities(context.TODO(), perPage, beforeTimePtr, afterTimePtr)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		activitiesJsonBytes, err := json.Marshal(activities)
		if err != nil {
			fmt.Println(err.Error())
			return
		}


		if outputPath != "" {
			utils.WriteOutput(outputPath, activitiesJsonBytes)
		}
		fmt.Println(string(activitiesJsonBytes))
	},
}
func init() {
	getActivities.Flags().IntVarP(&perPage, "per-page", "n", 30, "The number of activities to get per page. (max 200)")
	getActivities.Flags().StringVarP(&before, "before", "b", "", fmt.Sprintf("Filter to only include activities before this date. Must be of the format: %s", layoutInterpretation))
	getActivities.Flags().StringVarP(&after, "after", "a", "", fmt.Sprintf("Filter to only include activities after this date. Must be of the format: %s", layoutInterpretation))

	tokenCmdGroup.AddCommand(getActivities)
}