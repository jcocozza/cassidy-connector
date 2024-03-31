package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)
var numPages int
var perPage int
var outputPath string
// used to bulk grab activities
var getActivities = &cobra.Command{
	Use: "activity [access token]",
	Short: "Get activities. Expects an access token.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		accessToken := args[0]
		stravaApp := createApp()

		pages, err := stravaApp.GetActivityPages(accessToken, numPages, perPage)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var combined []byte
		for _, byteList := range pages {
			// Remove the enclosing square brackets from each byte slice
			byteList = byteList[1 : len(byteList)-1]
			// Append the byte slice to the combined list
			combined = append(combined, byteList...)
		}
		// Wrap the combined JSON objects within square brackets to form a single list
		combined = append([]byte{'['}, append(combined, ']')...)

		// Write the combined JSON data to a file
		err1 := os.WriteFile(outputPath, combined, 0644)
		if err1 != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Println(string(combined))
	},
}
func init() {
	getActivities.Flags().IntVarP(&numPages, "pages", "p", 1, "The number of pages to get.")
	getActivities.Flags().IntVarP(&perPage, "per-page", "n", 30, "The number of activities to get per page.")
	getActivities.Flags().StringVarP(&outputPath, "path", "f", "output.json", "The path to save the json output to.")

	rootCmd.AddCommand(getActivities)
}