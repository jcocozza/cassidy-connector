package cmd

import (
	"fmt"
	"os"

	stravaCmd "github.com/jcocozza/cassidy-connector/strava/cmd"
	finalSurgeCmd "github.com/jcocozza/cassidy-connector/finalSurge/cmd"
	"github.com/spf13/cobra"
)

const (
	version string = "0.0.1"
)

var rootCmd = &cobra.Command{
	Use: "cassidy",
	Version: version,
	Short: "cassidy is a cli tool for interacting with different activity API's",
	Long: `cassidy is a cli tool for interacting with different activity API's
The project is currently developing support for:
- Strava
- Final Surge`,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(stravaCmd.RootCmd)
	rootCmd.AddCommand(finalSurgeCmd.RootCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}