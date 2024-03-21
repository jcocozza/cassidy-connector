package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:   "cassidy-strava",
	Version: version,
	Short: "cassidy-strava is a cli tool to interact with the Strava API",
	Long: `cassidy-strava is a cli tool to interact with the Strava API`,
	Run: func(cmd *cobra.Command, args []string) {
	  // Do Stuff Here
	},
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
	  fmt.Fprintln(os.Stderr, err)
	  os.Exit(1)
	}
}