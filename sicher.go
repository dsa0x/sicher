/*
Sicher is a go module that allows safe storage of encrypted credentials in a version control system.
It is a port of the secret management system that was introduced in Ruby on Rails 6.
Examples can be found in examples/ folder
*/
package sicher

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
)

var delimiter = "==--=="
var defaultEnv = "dev"

type Sicher struct {
	// Path is the path to the project. Defaults to the current directory
	Path string

	// Environment is the environment to use. Defaults to "dev"
	Environment string
	data        map[string]string `yaml:"data"`
}

// New creates a new sicher struct
func New(environment string, path ...string) *Sicher {

	if environment == "" {
		environment = defaultEnv
	}

	var _path string
	if len(path) < 1 || path[0] == "" {
		_path = "."
	} else {
		_path = path[0]
	}
	_path, _ = filepath.Abs(_path)
	return &Sicher{Path: _path + "/", Environment: environment, data: make(map[string]string)}
}

// Initialize initializes the sicher project and creates the necessary files
func (s *Sicher) Initialize() {
	key := generateKey()

	// create the key file if it doesn't exist
	keyFile, err := os.OpenFile(fmt.Sprintf("%s%s.key", s.Path, s.Environment), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
		return
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
			log.Println(err)
			return
		}
		encFileFlag = os.O_CREATE | os.O_RDWR | os.O_TRUNC
	} else {
		// if keyfile exists, append to it
		encFileFlag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}

	// create the encrypted credentials file if it doesn't exist
	// the credentials file is truncated if the keyfile is new
	encFile, err := os.OpenFile(fmt.Sprintf("%s%s.enc", s.Path, s.Environment), encFileFlag, 0600)
	if err != nil {
		log.Println(err)
		return
	}
	defer encFile.Close()

	encFileStats, _ := encFile.Stat()

	// if the encrypted file is new, write some random data to it
	if encFileStats.Size() < 1 {
		initFile := []byte(`TESTKEY=loremipsum`)
		nonce, ciphertext, err := encrypt(key, initFile)
		if err != nil {
			fmt.Printf("Error encrypting file: %s\n", err)
			return
		}
		_, err = encFile.WriteString(fmt.Sprintf("%x%s%x", ciphertext, delimiter, nonce))
		if err != nil {
			log.Println(err)
			return
		}
	}

	// add the key file to gitignore
	f, err := os.OpenFile(fmt.Sprintf("%s.gitignore", s.Path), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
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
			log.Println(err)
			return
		}

		if string(str) == fmt.Sprintf("%s.key", s.Environment) {
			return
		}

		if err == io.EOF {
			break
		}
	}

	f.Write([]byte(fmt.Sprintf("\n%s.key", s.Environment)))
}

// Edit opens the encrypted credentials in a temporary file for editing. Default editor is vim.
func (s *Sicher) Edit(editor ...string) {
	var editorName string
	if len(editor) > 0 {
		editorName = editor[0]
	} else {
		editorName = "vim"
	}

	match, _ := regexp.MatchString("^(nano|vim|vi|)$", editorName)
	if !match {
		log.Println("Invalid Command: Select one of vim, vi, or nano as editor, or leave as empty")
		return
	}

	// read the encryption key
	key, err := os.ReadFile(fmt.Sprintf("%s.key", s.Environment))
	if err != nil {
		log.Printf("encryption key(%s.key) is not available. Create one by running the cli with init flag.\n", s.Environment)
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
	f, err := os.CreateTemp("", "*-credentials.env")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	filePath := f.Name()
	defer cleanUpFile(filePath)

	// if file already exists, decode and decrypt it
	nonce, fileText, err := decodeFile(enc)
	if err != nil {
		fmt.Printf("Error decoding encryption file: %s\n", err)
		return
	}

	if nonce != nil && fileText != nil {
		plaintext, err := decrypt(strKey, nonce, fileText)
		if err != nil {
			fmt.Println("Error decrypting file:", err)
			return
		}
		_, err = f.Write(plaintext)
		if err != nil {
			fmt.Printf("Error saving credentials: %s \n", err)
			return
		}
	}

	//open decrypted file with editor
	cmd := exec.Command(editorName, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error while editing %v \n", err)
		return
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	//encrypt and overwrite credentials file
	// the encrypted file is encoded in hexadecimal format
	nonce, encrypted, err := encrypt(strKey, file)
	if err != nil {
		fmt.Printf("Error encrypting file: %s \n", err)
		return
	}

	credFile.Truncate(0)
	credFile.Write([]byte(fmt.Sprintf("%x%s%x", encrypted, delimiter, nonce)))
	log.Println("File encrypted and saved.")

}

func (s *Sicher) loadEnv() {
	s.configure()
	s.setEnv()
}

// LoadEnv loads the environment variables from the encrypted credentials file into the config gile.
// configFile can be a struct or map[string]string
func (s *Sicher) LoadEnv(prefix string, configFile interface{}) error {
	s.loadEnv()

	d := reflect.ValueOf(configFile)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
	} else {
		return errors.New("configFile must be a pointer to a struct or map")
	}

	if !(d.Kind() == reflect.Struct || d.Kind() == reflect.Map) {
		return errors.New("config must be a type of struct or map")
	}

	// the configFile is a map, set the values
	if d.Kind() == reflect.Map {
		if d.Type() != reflect.TypeOf(map[string]string{}) {
			return errors.New("configFile must be a struct or map[string]string")
		}
		d.Set(reflect.ValueOf(s.data))
		return nil
	}

	// if the interface is a struct, iterate over the fields and set the values
	for i := 0; i < d.NumField(); i++ {
		field := d.Field(i)
		fieldType := d.Type().Field(i)
		isRequired := fieldType.Tag.Get("required")
		key := fieldType.Tag.Get("envconfig")

		tagName := key
		if prefix != "" {
			tagName = fmt.Sprintf("%s_%s", prefix, key)
		}

		envVar := os.Getenv(tagName)
		if isRequired == "true" && envVar == "" {
			return errors.New("required env variable " + key + " is not set")
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envVar)
		case reflect.Bool:
			field.SetBool(envVar == "true")
		}

	}
	return nil
}
