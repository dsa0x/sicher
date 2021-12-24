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
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"time"

	"github.com/juju/fslock"
)

var delimiter = "==--=="
var defaultEnv = "dev"
var DefaultEnvStyle = DOTENV
var masterKey = "SICHER_MASTER_KEY"
var (
	execCmd               = exec.Command
	stdIn   io.ReadWriter = os.Stdin
	stdOut  io.ReadWriter = os.Stdout
	stdErr  io.ReadWriter = os.Stderr
)

type sicher struct {
	// Path is the path to the project. If empty string, it defaults to the current directory
	Path string

	// Environment is the environment to use. Defaults to "dev"
	Environment string
	data        map[string]string `yaml:"data"`

	envStyle EnvStyle

	// gitignorePath is the path to the .gitignore file
	gitignorePath string
}

// New creates a new sicher struct
// path is the path to the project. If empty string, it defaults to the current directory
// environment is the environment to use. Defaults to "dev"
func New(environment string, path string) *sicher {

	if environment == "" {
		environment = defaultEnv
	}

	if path == "" {
		path = "."
	}
	path, _ = filepath.Abs(path)
	return &sicher{Path: path + "/", Environment: environment, data: make(map[string]string), envStyle: DOTENV}
}

// Initialize initializes the sicher project and creates the necessary files
func (s *sicher) Initialize(scanReader io.Reader) error {
	key := generateKey()

	// create the key file if it doesn't exist
	keyFile, err := os.OpenFile(fmt.Sprintf("%s%s.key", s.Path, s.Environment), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("error creating key file: %s", err)
	}
	defer keyFile.Close()

	keyFileStats, err := keyFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting key file stats: %s", err)
	}

	// create the encrypted credentials file if it doesn't exist
	encFile, err := os.OpenFile(fmt.Sprintf("%s%s.enc", s.Path, s.Environment), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("error creating encrypted credentials file: %s", err)
	}
	defer encFile.Close()

	encFileStats, err := encFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting key file stats: %s", err)
	}

	// if keyfile is new
	// Absence of keyfile indicates that the project is new or keyfile is lost
	// if keyfile is lost, the encrypted file cannot be decrypted,
	// and the user needs to re-initialize or obtain the original key
	if keyFileStats.Size() < 1 {

		// if encrypted file exists
		// ask user if they want to overwrite the encrypted file
		// if yes, truncate file and continue
		// else cancel
		if encFileStats.Size() > 1 {
			fmt.Printf("An encrypted credentials file already exist, do you want to overwrite it? \n Enter 'yes' or 'y' to accept.\n")
			rd := bufio.NewScanner(scanReader)
			for rd.Scan() {
				line := rd.Text()
				if line == "yes" || line == "y" {
					encFile.Truncate(0)
					break
				} else {
					os.Remove(keyFile.Name())
					fmt.Println("Exiting. Leaving credentials file unmodified")
					return nil
				}
			}
		}

		_, err = keyFile.WriteString(key)
		if err != nil {
			return fmt.Errorf("error saving key file: %s", err)
		}
	}

	// stats will have changed if the file was truncated
	encFileStats, err = encFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting key file stats: %s", err)
	}

	// if the encrypted file is new, write some random data to it
	if encFileStats.Size() < 1 {
		initFile := []byte(fmt.Sprintf("TESTKEY%sloremipsum\n", envStyleDelim[s.envStyle]))
		nonce, ciphertext, err := encrypt(key, initFile)
		if err != nil {
			return fmt.Errorf("error encrypting credentials file: %s", err)
		}
		_, err = encFile.WriteString(fmt.Sprintf("%x%s%x", ciphertext, delimiter, nonce))
		if err != nil {

			return fmt.Errorf("error writing encrypted credentials file: %s", err)
		}
	}

	// add the key file to gitignore
	if s.gitignorePath != "" {
		err = addToGitignore(fmt.Sprintf("%s.key", s.Environment), s.gitignorePath)
		if err != nil {
			return fmt.Errorf("error adding key file to gitignore: %s", err)
		}
	}
	return nil
}

// Edit opens the encrypted credentials in a temporary file for editing. Default editor is vim.
func (s *sicher) Edit(editor ...string) error {
	var editorName string
	if len(editor) > 0 {
		editorName = editor[0]
	} else {
		editorName = "vim"
	}

	match, _ := regexp.MatchString("^(nano|vim|vi|code|)$", editorName)
	if !match {
		return fmt.Errorf("invalid Command: Select one of vim, vi, code or nano as editor, or leave as empty")
	}

	var cmdArgs []string
	// waitOpt is needed to enable vscode to wait for the editor to close before continuing
	var waitOpt string
	if editorName == "code" {
		waitOpt = "--wait"
		cmdArgs = append(cmdArgs, waitOpt)
	}

	// read the encryption key. if key not in file, try getting from env
	key, err := os.ReadFile(fmt.Sprintf("%s%s.key", s.Path, s.Environment))
	if err != nil {
		if os.Getenv(masterKey) != "" {
			key = []byte(os.Getenv(masterKey))
		} else {
			return fmt.Errorf("encryption key(%s.key) is not available. Provide a key file or enter one through the command line", s.Environment)
		}
	}
	strKey := string(key)

	// open the encrypted credentials file
	credFile, err := os.OpenFile(fmt.Sprintf("%s%s.enc", s.Path, s.Environment), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer credFile.Close()

	// lock file to enable only one edit at a time
	credFileLock := fslock.New(credFile.Name()) //newFileLock(credFile)
	err = credFileLock.LockWithTimeout(time.Second * 10)
	if err != nil {
		return fmt.Errorf("error locking file: %s", err)
	}
	defer credFileLock.Unlock()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, credFile)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	enc := buf.String()

	// Create a temporary file to edit the decrypted credentials
	f, err := os.CreateTemp("", fmt.Sprintf("*-credentials.%s", envStyleExt[s.envStyle]))
	if err != nil {
		return fmt.Errorf("error creating temp file %v", err)
	}
	defer f.Close()
	filePath := f.Name()
	defer cleanUpFile(filePath)

	// if file already exists, decode and decrypt it
	nonce, fileText, err := decodeFile(enc)
	if err != nil {
		return fmt.Errorf("error decoding encryption file: %s", err)
	}

	var plaintext []byte
	if nonce != nil && fileText != nil {
		plaintext, err = decrypt(strKey, nonce, fileText)
		if err != nil {
			return fmt.Errorf("error decrypting file: %s", err)
		}

		_, err = f.Write(plaintext)
		if err != nil {
			return fmt.Errorf("error saving credentials: %s", err)
		}
	}

	//open decrypted file with editor
	cmdArgs = append(cmdArgs, filePath)
	cmd := execCmd(editorName, cmdArgs...)
	cmd.Stdin = stdIn
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting editor: %s", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error while editing %v", err)
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading credentials file %v ", err)
	}

	// if no file changes, dont generate new encrypted file
	if bytes.Equal(file, plaintext) {
		fmt.Fprintf(stdOut, "No changes made.\n")
		return nil
	}

	//encrypt and overwrite credentials file
	// the encrypted file is encoded in hexadecimal format
	nonce, encrypted, err := encrypt(strKey, file)
	if err != nil {
		return fmt.Errorf("error encrypting file: %s ", err)
	}

	credFile.Truncate(0)
	credFile.Write([]byte(fmt.Sprintf("%x%s%x", encrypted, delimiter, nonce)))
	fmt.Fprintf(stdOut, "File encrypted and saved.\n")
	return nil
}

// LoadEnv loads the environment variables from the encrypted credentials file into the config gile.
// configFile can be a struct or map[string]string
func (s *sicher) LoadEnv(prefix string, configFile interface{}) error {
	s.configure()
	s.setEnv()

	d := reflect.ValueOf(configFile)
	if d.Kind() == reflect.Ptr {
		d = d.Elem()
	} else {
		return errors.New("configFile must be a pointer to a struct or map")
	}

	if !(d.Kind() == reflect.Struct || d.Kind() == reflect.Map) {
		return errors.New("config must be a type of struct or map")
	}

	// the configFile is a map, set the values and return
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
		key := fieldType.Tag.Get("env")

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

func (s *sicher) SetEnvStyle(style string) {
	if style != "dotenv" && style != "yaml" && style != "yml" {
		fmt.Println("Invalid style: Select one of dotenv, yml, or yaml")
		os.Exit(1)
	}
	s.envStyle = EnvStyle(style)
}

func (s *sicher) SetGitignorePath(path string) {
	path, _ = filepath.Abs(path)
	s.gitignorePath = path
}
