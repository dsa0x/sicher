package sicher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// encrypt encrypts the given plaintext with the given key and returns the ciphertext
func encrypt(key string, fileData []byte) (nonce []byte, ciphertext []byte, err error) {
	hKey, err := hex.DecodeString(key)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(hKey)
	if err != nil {
		return
	}

	nonce = make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	ciphertext = aesgcm.Seal(nil, nonce, fileData, nil)
	return
}

func decrypt(key string, nonce, text []byte) (plaintext []byte, err error) {
	hKey, err := hex.DecodeString(key)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(hKey)
	if err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	plaintext, err = aesgcm.Open(nil, nonce, text, nil)
	return
}
