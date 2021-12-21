package sicher

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// configure reads the credentials file and sets the environment variables
func (s *Sicher) configure() {

	if s.Environment == "" {
		s.Environment = "development"
	}
	// read the encryption key
	key, err := os.ReadFile(fmt.Sprintf("%s.key", s.Environment))
	if err != nil {
		fmt.Printf("Encryption key (%s.key) is not available. Create one by running the cli with init flag.", s.Environment)
		return
	}
	strKey := string(key)

	// read the encrypted credentials file
	credFile, err := os.ReadFile(fmt.Sprintf("%s.enc", s.Environment))
	if err != nil {
		fmt.Printf("Encrypted credentials file (%s.enc) is not available. Create one by running the cli with init flag.", s.Environment)
		return
	}

	encFile := string(credFile)

	var envBuf bytes.Buffer

	// if file already exists, decode and decrypt it
	nonce, fileText, err := decodeFile(encFile)
	if err != nil {
		fmt.Printf("Error decoding encryption file: %s\n", err)
		return
	}

	if nonce != nil && fileText != nil {
		plaintext := Decrypt(strKey, nonce, fileText)
		_, err = envBuf.Write(plaintext)
		if err != nil {
			fmt.Printf("Error decoding credentials: %s\n", err)
			return
		}
	}

	err = yaml.Unmarshal(envBuf.Bytes(), &s.data)
	if err != nil {
		fmt.Printf("Error decoding credentials: %s\n", err)
	}
}

func (s *Sicher) setEnv() {
	for k, v := range s.data {
		err := os.Setenv(k, fmt.Sprintf("%v", v))
		if err != nil {
			panic(err)
		}
	}
}
