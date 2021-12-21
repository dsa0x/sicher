package cli

import (
	"fmt"
	"os"

	"github.com/dsa0x/sicher/sicher"
	"github.com/spf13/cobra"
)

var sich = sicher.New()

var rootCmd = &cobra.Command{
	Use:   "sicher",
	Short: "Sicher is a tool for managing your Go projects",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	sich.SetCredentials()
	fmt.Println(os.Getenv("base_key"))
}

func testEnv() {
	fmt.Println(os.Getenv("mongodb"), "ENV")
}
