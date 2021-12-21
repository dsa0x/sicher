package sicher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
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
