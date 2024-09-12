package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	version string = "0.0.1"
)

var outputPath string

var email string
var password string

var RootCmd = &cobra.Command{
	Use: "cassidy-final-surge",
	Short: "cassidy-final-surge is a cli tool for interacting with the final surge API",
	Long: "warning: This is a 'back-engineered' tool. It can break at any time for pretty much any reason because Final Surge does not expose any standard procedures for users",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&email, "email", "", "your final surge email")
	RootCmd.PersistentFlags().StringVar(&password, "password", "", "your final surge password")

	RootCmd.PersistentFlags().StringVarP(&outputPath, "path", "f", "", "the path to save successful output to. (will not write errors at this time)")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
