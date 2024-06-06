package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jcocozza/cassidy-connector/strava/app/api"
	"github.com/jcocozza/cassidy-connector/utils"
	"github.com/spf13/cobra"
)

var keys []string
var getStreams = &cobra.Command{
	Use: "streams [activity id]",
	Short: "Get streams for a given activity",
	Long: `Get streams for a given activity. Stream types include:
time, distance, latlng, altitude, velocity_smooth, heartrate, cadence, watts, temp, moving, and grade_smooth`,
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

		l := []api.StreamType{}
		for _, key := range keys {
			l = append(l, api.StreamType(key))
		}
		streams, err := stravaApp.Api.GetActivityStreams(context.TODO(), activityId, l)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		streamJsonBytes, err := json.Marshal(streams)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if outputPath != "" {
			utils.WriteOutput(outputPath, streamJsonBytes)
		}
		fmt.Println(string(streamJsonBytes))
	},
}
func init() {
	getStreams.Flags().StringSliceVarP(&keys, "stream-types", "t", []string{"time", "distance"}, "a comma separated list of the stream types to get.")

	tokenCmdGroup.AddCommand(getStreams)
}