package sicher

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	data map[string]interface{} `yaml:"data"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Get() {
	// read the encryption key
	key, err := os.ReadFile(fmt.Sprintf("%s.key", environment))
	if err != nil {
		log.Fatal(err)
	}
	strKey := string(key)

	// read the encrypted credentials file
	credFile, err := os.ReadFile(fmt.Sprintf("%s.enc", environment))
	if err != nil {
		log.Fatalln(err)
	}

	encFile := string(credFile)

	var envBuf bytes.Buffer

	err = decodeFile(strKey, encFile, &envBuf)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = yaml.Unmarshal(envBuf.Bytes(), &c.data)
	if err != nil {
		panic(err)
	}
}

func (c *Config) SetEnv() {
	for k, v := range c.data {
		err := os.Setenv(k, fmt.Sprintf("%v", v))
		if err != nil {
			panic(err)
		}
	}
}
