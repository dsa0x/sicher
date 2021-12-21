package cli

import (
	"github.com/spf13/cobra"
)

var environment string
var path string

func init() {

	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&environment, "env", "development", "Enter your deployment environment")
	initCmd.Flags().StringVar(&path, "path", ".", "Enter the path to your project")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sicher in your project",
	Long:  `Initialize a new project`,
	Run: func(cmd *cobra.Command, args []string) {
		sich.Environment = environment
		sich.Initialize(path)
	},
}
