package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)


const layout string = "2006-01-02"
const layoutInterpretation string = "YYYY-MM-DD"

var start string
var end string
var getActivities = &cobra.Command{
	Use: "activities",
	Short: "get user activities",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		finalSurgeApp := createApp(email, password)
		auth, err := finalSurgeApp.Authenticate(context.TODO())
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		startDate, err := time.Parse(layout, start)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		endDate, err := time.Parse(layout, end)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		activities, err := finalSurgeApp.GetActivities(context.TODO(), auth.Data.Token, auth.Data.UserKey, startDate, endDate)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(string(activities))
	},
}

func init() {
	rootCmd.AddCommand(getActivities)

	getActivities.Flags().StringVarP(&start, "start", "s", "", fmt.Sprintf("Filter to only include activities after this date. Must be of the format: %s", layoutInterpretation))
	getActivities.Flags().StringVarP(&end, "end", "e", "", fmt.Sprintf("Filter to only include activities before this date. Must be of the format: %s", layoutInterpretation))

	getActivities.MarkFlagRequired("start")
	getActivities.MarkFlagRequired("end")
}