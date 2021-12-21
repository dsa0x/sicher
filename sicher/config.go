package sicher

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// configure reads the credentials file and sets the environment variables
func (s *Sicher) configure() {
	// read the encryption key
	key, err := os.ReadFile(fmt.Sprintf("%s.key", s.Environment))
	if err != nil {
		log.Fatal(err)
	}
	strKey := string(key)

	// read the encrypted credentials file
	credFile, err := os.ReadFile(fmt.Sprintf("%s.enc", s.Environment))
	if err != nil {
		log.Fatalln(err)
	}

	encFile := string(credFile)

	var envBuf bytes.Buffer

	// if file already exists, decode and decrypt it
	nonce, fileText, err := decodeFile(encFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if nonce != nil && fileText != nil {
		plaintext := Decrypt(strKey, nonce, fileText)
		_, err = envBuf.Write(plaintext)
		if err != nil {
			log.Printf("Error saving credentials: %s", err)
			return
		}
	}

	err = yaml.Unmarshal(envBuf.Bytes(), &s.data)
	if err != nil {
		panic(err)
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
