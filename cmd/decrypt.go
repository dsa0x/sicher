package sicher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	data map[string]interface{} `yaml:"data"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Get() {
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err.Error())
	}

	yaml.Unmarshal(file, &c.data)
	fmt.Println(c.data)
	err = yaml.Unmarshal(file, c)
	if err != nil {
		panic(err.Error())
	}
}

func (c *Config) toString() {

}

func (c *Config) Encrypt(key string) (nonce []byte, ciphertext []byte) {
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

	data, err := json.Marshal(c.data)
	if err != nil {
		log.Printf("Marshal Error: %s", err.Error())
		return
	}

	ciphertext = aesgcm.Seal(nil, nonce, []byte(data), nil)
	return
}

func (c *Config) Decrypt(key string, nonce []byte, text []byte) (plaintext []byte) {
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
