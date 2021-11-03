package sicher

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type cred map[string]string

var stages = map[string]cred{
	"dev": {
		"keyName":        "development.key",
		"credentialName": "development.enc",
	},
	"prod": {
		"keyName":        "development.key",
		"credentialName": "development.enc",
	},
	"staging": {
		"keyName":        "development.key",
		"credentialName": "development.enc",
	},
}
var environment string

func init() {
	// cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&environment, "env", "development", "Enter your deployment environment")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sicher in your project",
	Long:  `Initialize a new project`,
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	key := generateKey()

	err := ioutil.WriteFile(fmt.Sprintf("%s.key", environment), []byte(key), 0600)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s.enc", environment), []byte(""), 0600)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	f.Write([]byte(fmt.Sprintf("\n%s.key\n", environment)))
}

func generateKey() string {
	timestamp := time.Now().UnixNano()
	key := sha256.Sum256([]byte(fmt.Sprint(timestamp)))
	return fmt.Sprintf("%x", key)
}
