package sicher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

type Encryption struct {
	CipherType string
	CipherKey  string
}

func NewEncryption(cipherType string) *Encryption {

	if cipherType == "" {
		cipherType = "aes-256-gcm"
	}

	return &Encryption{
		CipherType: cipherType,
	}
}

func GenerateKey() string {
	timestamp := time.Now().UnixNano()
	key := sha256.Sum256([]byte(fmt.Sprint(timestamp)))
	return fmt.Sprintf("%x", key)
}

func Encrypt(key string, fileData []byte) (nonce []byte, ciphertext []byte) {
	hKey, _ := hex.DecodeString(key)
	block, err := aes.NewCipher(hKey)
	if err != nil {
		panic(err.Error())
	}

	nonce = make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	// data, err := json.Marshal(fileData)
	// if err != nil {
	// 	log.Printf("Marshal Error: %s", err.Error())
	// 	return
	// }

	ciphertext = aesgcm.Seal(nil, nonce, fileData, nil)
	return
}

func Decrypt(key string, nonce []byte, text []byte) (plaintext []byte) {
	hKey, _ := hex.DecodeString(key)
	block, err := aes.NewCipher(hKey)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err = aesgcm.Open(nil, nonce, text, nil)
	if err != nil {
		panic(err.Error())
	}
	return
}
