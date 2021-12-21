package cli

import (
	"github.com/spf13/cobra"
)

var editor string

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringVar(&environment, "env", "development", "Enter your deployment environment")
	editCmd.Flags().StringVar(&editor, "editor", "vim", "Select editor. vim | vi | nano")
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit credentials",
	Run: func(cmd *cobra.Command, args []string) {
		sich.Environment = environment
		sich.Edit(editor)
	},
}
