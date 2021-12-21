package cli

import (
	"fmt"
	"os"

	"github.com/dsa0x/sicher/sicher"
	"github.com/spf13/cobra"
)

// use default values
var sich = &sicher.Sicher{Environment: "development", Path: "."}

var rootCmd = &cobra.Command{
	Use:   "sicher",
	Short: "sicher is a tool for encrypting and managing environment variables",
}

func init() {
	rootCmd.Example = `
	# Initialize sicher in your project
	sicher init --env development --path .

	# Edit environment variables
	sicher edit --env development --path . --editor vim
	`
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func LoadEnv() {
	sich.LoadEnv()
}
