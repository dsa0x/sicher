package sicher

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sicher",
	Short: "Sicher is a tool for managing your Go projects",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type Sicher struct {
	path        string
	Environment string
	data        map[string]interface{} `yaml:"data"`
}

func New() *Sicher {
	return &Sicher{}
}

// Initialize initializes the sicher project and creates the necessary files
func (s *Sicher) Initialize(path ...string) {
	key := generateKey()

	if len(path) > 0 {
		dir, _ := filepath.Abs(path[0])
		s.path = dir + "/"
	}

	if s.Environment == "" {
		s.Environment = "development"
	}
	// create the key file if it doesn't exist
	keyFile, err := os.OpenFile(fmt.Sprintf("%s%s.key", s.path, s.Environment), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer keyFile.Close()

	keyFileStats, _ := keyFile.Stat()
	encFileFlag := 0

	// if keyfile is new, write the key and truncate the encrypted file so it can be freshly written to
	// Absence of keyfile indicates that the project is new or keyfile is lost
	// if keyfile is lost, the encrypted file cannot be decrypted and the user needs to re-initialize
	if keyFileStats.Size() < 1 {
		_, err = keyFile.WriteString(key)
		if err != nil {
			log.Fatal(err)
		}
		encFileFlag = os.O_CREATE | os.O_RDWR | os.O_TRUNC
	} else {
		// if keyfile exists, append to it
		encFileFlag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}

	// create the encrypted credentials file if it doesn't exist
	// the credentials file is truncated if the keyfile is new
	encFile, err := os.OpenFile(fmt.Sprintf("%s%s.enc", s.path, s.Environment), encFileFlag, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer encFile.Close()

	encFileStats, _ := encFile.Stat()

	// if the encrypted file is new, write some random data to it
	if encFileStats.Size() < 1 {
		initFile := []byte(`base_key: test key`)
		nonce, ciphertext := Encrypt(key, initFile)
		_, err = encFile.WriteString(fmt.Sprintf("%x\n%x", ciphertext, nonce))
		if err != nil {
			log.Fatal(err)
		}
	}

	// add the key file to gitignore
	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	fr := bufio.NewReader(f)

	// check if the key file is already in the .gitignore file before adding it
	// if it is, don't add it again
	for err == nil {
		str, _, err := fr.ReadLine()
		if err != nil && err != io.EOF {
			log.Fatalln(err)
			return
		}

		if string(str) == fmt.Sprintf("%s%s.key", s.path, s.Environment) {
			return
		}

		if err == io.EOF {
			break
		}
	}

	f.Write([]byte(fmt.Sprintf("\n%s%s.key", s.path, s.Environment)))
}

// Edit opens the encrypted credentials file. Default editor is vim.
func (s *Sicher) Edit(editor ...string) {
	var editorName string
	if len(editor) > 0 {
		editorName = editor[0]
	} else {
		editorName = "vim"
	}

	if s.Environment == "" {
		s.Environment = "development"
	}

	match, _ := regexp.MatchString("^(nano|vim|vi|)$", editorName)
	if !match {
		log.Println("Invalid Command: Select one of vim, vi, or nano as editor, or leave as empty")
		return
	}

	// read the encryption key
	key, err := os.ReadFile(fmt.Sprintf("%s.key", s.Environment))
	if err != nil {
		log.Printf("Encryption key (%s.key) is not available. Create one by running the cli with init flag.", s.Environment)
		return
	}
	strKey := string(key)

	// open the encrypted credentials file
	credFile, err := os.OpenFile(fmt.Sprintf("%s.enc", s.Environment), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
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
	nonce, fileText, err := decodeFile(enc)
	if err != nil {
		fmt.Println(err)
		return
	}

	if nonce != nil && fileText != nil {
		plaintext := Decrypt(strKey, nonce, fileText)
		_, err = f.Write(plaintext)
		if err != nil {
			log.Printf("Error saving credentials: %s", err)
			return
		}
	}

	//open decrypted file with editor
	command := exec.Command(editorName, filePath)
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
	log.Println("File encrypted and saved")

}

func (s *Sicher) SetCredentials() {
	s.configure()
	s.setEnv()
}
