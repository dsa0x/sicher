package sicher

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var editor string

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringVar(&environment, "env", "development", "Enter your deployment environment")
	editCmd.Flags().StringVar(&editor, "editor", "vim", "Select editor. vim | vi | nano")
}

var editCmd = &cobra.Command{
	Use: "edit",
	// Usage: "edit [OPTIONS] [FILE]",
	Short: "Edit credentials",
	Run: func(cmd *cobra.Command, args []string) {
		Edit(args)
	},
}

func Edit(args []string) {

	match, _ := regexp.MatchString("^(nano|vim|vi|)$", editor)
	if !match {
		log.Println("Invalid Command: Select one of vim, vi, or nano as editor, or leave as empty")
		return
	}

	// read the encryption key
	key, err := os.ReadFile(fmt.Sprintf("%s.key", environment))
	if err != nil {
		log.Fatal(err)
	}
	strKey := string(key)

	// open the encrypted credentials file
	credFile, err := os.OpenFile(fmt.Sprintf("%s.enc", environment), os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer credFile.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, credFile)
	if err != nil {
		log.Fatalln(err)
	}
	enc := buf.String()

	// Create a temporary file to edit the decrypted credentials
	f, err := os.CreateTemp("", "*-credentials.yml")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	filePath := f.Name()
	defer cleanUpFile(filePath)

	// if file already exists, decode and decrypt it
	err = decodeFile(strKey, enc, f)
	if err != nil {
		fmt.Println(err)
		return
	}

	//open decrypted file with editor
	command := exec.Command(editor, filePath)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Start()
	if err != nil {
		log.Println(err)
		return
	}

	err = command.Wait()
	if err != nil {
		log.Printf("Error while editing %v", err)
		return
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	//encrypt and overwrite credentials file
	nonce, encrypted := Encrypt(strKey, file)
	str := hex.EncodeToString(encrypted)
	credFile.Truncate(0)
	credFile.Write([]byte(fmt.Sprintf("%s\n%s", str, hex.EncodeToString(nonce))))
	log.Printf("File encrypted and saved")

}

func generateRandomPath() string {
	b := make([]byte, 16)
	rand.Read(b)
	return os.TempDir() + base64.RawURLEncoding.EncodeToString(b)[:] + "-credentials.yml"
}

func cleanUpFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Error while cleaning up %v", err.Error())
	}
}

func decodeFile(encKey, encFile string, f io.Writer) error {
	if len(encFile) > 0 {
		resp := strings.Split(encFile, "\n")
		if len(resp) < 2 {
			return errors.New("Invalid credentials")
		}
		nonce, err := hex.DecodeString(resp[1])
		if err != nil {
			log.Printf("Invalid credentials file: %s", err)
			return err
		}
		res, err := hex.DecodeString(resp[0])
		if err != nil {
			log.Printf("Invalid credentials file: %s", err)
			return err
		}
		plaintext := Decrypt(encKey, nonce, res)
		_, err = f.Write(plaintext)
		if err != nil {
			log.Printf("Error saving credentials: %s", err)
			return err
		}
	}
	return nil
}
