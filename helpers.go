package sicher

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type EnvStyle string

const (
	YAML  EnvStyle = "yaml"
	YML   EnvStyle = "yml"
	BASIC EnvStyle = "basic"
)

var envStyleDelim = map[EnvStyle]string{
	YAML:  ":",
	YML:   ":",
	BASIC: "=",
}

var envStyleExt = map[EnvStyle]string{
	YAML:  "yml",
	YML:   "yml",
	BASIC: "env",
}

// cleanUpFile removes the given file
func cleanUpFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Printf("Error while cleaning up %v \n", err.Error())
	}
}

// decodeFile decodes the encrypted file and returns the decoded file and nonce
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

// generateKey generates a random key of 32 bytes and encodes as hex string
func generateKey() string {
	timestamp := time.Now().UnixNano()
	key := sha256.Sum256([]byte(fmt.Sprint(timestamp)))
	rand.Read(key[16:])
	return hex.EncodeToString(key[:])
}

// parseConfig parses the environment variables into a map
func parseConfig(config []byte, store map[string]string, envType EnvStyle) (err error) {

	if envType != BASIC && envType != YAML && envType != YML {
		return errors.New("invalid environment type")
	}

	delim := envStyleDelim[envType]

	var b bytes.Buffer
	b.Write(config)
	sc := bufio.NewScanner(&b)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		cfgLine := strings.Split(line, delim)

		// ignore commented lines and invalid lines
		if len(cfgLine) < 2 || canIgnore(line) {
			continue
		}

		store[cfgLine[0]] = strings.Join(cfgLine[1:], delim)
		if err == io.EOF {
			return nil
		}
	}
	return nil

}

// canIgnore ignores commented lines and empty lines
func canIgnore(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, `#`) || len(line) == 0
}
