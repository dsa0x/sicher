package sicher

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
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
	editCmd.Flags().StringVar(&editor, "editor", "nano", "Select editor. vim | vi | nano")
}

var editCmd = &cobra.Command{
	Use: "edit",
	// Usage: "edit [OPTIONS] [FILE]",
	Short: "Edit credentials",
	Run:   runEdit,
}

func runEdit(cmd *cobra.Command, args []string) {

	match, _ := regexp.MatchString("^(nano|vim|vi|)$", editor)
	if !match {
		log.Println("Invalid Command: Select one of vim, vi, or nano as editor, or leave as empty")
		return
	}

	key, err := ioutil.ReadFile(fmt.Sprintf("%s.key", environment))
	if err != nil {
		log.Fatal(err)
	}
	strKey := string(key)

	enc, err := ioutil.ReadFile(fmt.Sprintf("%s.enc", environment))
	if err != nil {
		log.Fatalln(err)
	}
	credFile, err := os.OpenFile(fmt.Sprintf("%s.enc", environment), os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer credFile.Close()

	filePath := generateRandomPath()
	defer cleanUpFile(filePath)

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	if len(enc) > 0 {
		resp := strings.Split(string(enc), "\n")
		if len(resp) < 2 {
			log.Fatalln("Invalid credentials")
			return
		}
		nonce, _ := hex.DecodeString(resp[1])
		res, _ := hex.DecodeString(resp[0])
		plaintext := Decrypt(strKey, nonce, res)
		_, err = f.Write(plaintext)
		if err != nil {
			log.Fatal(err)
		}
	}

	//open file with editor
	command := exec.Command(editor, filePath)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Start()
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = command.Wait()
	if err != nil {
		log.Printf("Error while editing %v", err.Error())
		return
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	//encrypt and overwrite file
	nonce, encrypted := Encrypt(strKey, file)
	str := hex.EncodeToString(encrypted)
	credFile.Truncate(0)
	credFile.Write([]byte(fmt.Sprintf("%s\n%s", str, hex.EncodeToString(nonce))))
	log.Printf("File encrypted and saved")

}

func generateRandomPath() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return os.TempDir() + base64.RawURLEncoding.EncodeToString(b)[:8] + "-credentials.yml"
}

func cleanUpFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Error while cleaning up %v", err.Error())
	}
}
