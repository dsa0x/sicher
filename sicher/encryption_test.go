package sicher

import (
	"bytes"
	"testing"
)

func TestEncryption(t *testing.T) {

	key := generateKey()
	fileText := []byte("mytestfiletext")
	nonce, cipherText, err := encrypt(key, fileText)
	if err != nil {
		t.Errorf("Unable to encrypt file; got error %v", err)
	}

	plaintext, err := decrypt(key, nonce, cipherText)
	if err != nil {
		t.Errorf("Unable to decrypt file; got error %v", err)
	}

	if !bytes.Equal(fileText, plaintext) {
		t.Errorf("Expected fileText to be equal to plaintext, got %s and %s", fileText, plaintext)
	}

	// decrypting with an incorrect key
	_, err = decrypt(generateKey(), nonce, cipherText)
	if err == nil {
		t.Errorf("Expected ciphertext not to be decryptable using an incorrect key")
	}

}
