package sicher

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// cleanUpFile removes the given file
func cleanUpFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Printf("Error while cleaning up %v \n", err.Error())
	}
}

// decodeFile decodes the encrypted file and writes it to the given writer
func decodeFile(encFile string) (nonce []byte, fileText []byte, err error) {
	if encFile == "" {
		return nil, nil, nil
	}

	resp := strings.Split(encFile, delimiter)
	if len(resp) < 2 {
		return nil, nil, errors.New("invalid credentials")
	}
	nonce, err = hex.DecodeString(resp[1])
	if err != nil {
		return nil, nil, err
	}
	fileText, err = hex.DecodeString(resp[0])
	if err != nil {
		return nil, nil, err
	}

	return
}

func generateKey() string {
	timestamp := time.Now().UnixNano()
	key := sha256.Sum256([]byte(fmt.Sprint(timestamp)))
	rand.Read(key[16:])
	return fmt.Sprintf("%x", key)
}
