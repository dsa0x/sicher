package sicher

import (
	"fmt"
	"log"
	"os"
)

// configure reads the credentials file and sets the environment variables
func (s *sicher) configure() {

	if s.Environment == "" {
		fmt.Println("Environment not set")
		return
	}
	// read the encryption key
	fmt.Printf("%s%s.key", s.Path, s.Environment)
	key, err := os.ReadFile(fmt.Sprintf("%s%s.key", s.Path, s.Environment))
	if err != nil {
		fmt.Printf("encryption key (%s.key) is not available. Create one by running the cli with init flag.\n", s.Environment)
		return
	}
	strKey := string(key)

	// read the encrypted credentials file
	credFile, err := os.ReadFile(fmt.Sprintf("%s%s.enc", s.Path, s.Environment))
	if err != nil {
		fmt.Printf("encrypted credentials file (%s.enc) is not available. Create one by running the cli with init flag.\n", s.Environment)
		return
	}

	encFile := string(credFile)

	// if file already exists, decode and decrypt it
	nonce, fileText, err := decodeFile(encFile)
	if err != nil {
		fmt.Printf("Error decoding encryption file: %s\n", err)
		return
	}

	if nonce == nil || fileText == nil {
		fmt.Println("Error decoding encryption file: encrypted file is invalid")
		return
	}

	plaintext, err := decrypt(strKey, nonce, fileText)
	if err != nil {
		fmt.Println("Error decrypting file:", err)
		return
	}

	err = parseConfig(plaintext, s.data)
	if err != nil {
		fmt.Printf("Error decoding credentials: %s\n", err)
		return
	}

}

func (s *sicher) setEnv() {
	for k, v := range s.data {
		err := os.Setenv(k, fmt.Sprintf("%v", v))
		if err != nil {
			log.Fatalf("Error setting environment variable key %s: %s\n", k, err)
		}
	}
}
