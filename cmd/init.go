package sicher

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var environment string

func init() {

	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&environment, "env", "development", "Enter your deployment environment")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sicher in your project",
	Long:  `Initialize a new project`,
	Run: func(cmd *cobra.Command, args []string) {
		Initialize(args)
	},
}

func Initialize(args []string) {
	key := generateKey()

	// create the key file if it doesn't exist
	keyFile, err := os.OpenFile(fmt.Sprintf("%s.key", environment), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer keyFile.Close()

	keyFileStats, _ := keyFile.Stat()

	if keyFileStats.Size() < 1 {
		_, err = keyFile.WriteString(key)
		if err != nil {
			log.Fatal(err)
		}
	}

	// create the encrypted credentials file if it doesn't exist
	_, err = os.OpenFile(fmt.Sprintf("%s.enc", environment), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	// add the key file to gitignore
	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	f.Write([]byte(fmt.Sprintf("\n%s.key", environment)))
}

func generateKey() string {
	timestamp := time.Now().UnixNano()
	key := sha256.Sum256([]byte(fmt.Sprint(timestamp)))
	rand.Read(key[16:])
	return fmt.Sprintf("%x", key)
}
