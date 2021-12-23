package sicher

import (
	"encoding/hex"
	"os"
	"testing"
)

func TestCleanUpFile(t *testing.T) {
	f, err := os.CreateTemp("", "*tempfile.env")
	if err != nil {
		t.Errorf("Unable to create temporary test file; %v", err)
	}
	_, err = f.WriteString("Hello World")
	if err != nil {
		t.Errorf("Unable to write to temporary test file; %v", err)
	}

	// clean up test, not to be mistaken with the cleanup file function
	t.Cleanup(func() {
		f.Close()
	})

	cleanUpFile(f.Name())

	_, err = os.Open(f.Name())
	if err == nil {
		t.Errorf("file cleanup unsuccesssful")
	}

}

func TestGenerateKey(t *testing.T) {
	key := generateKey()
	_, err := hex.DecodeString(key)
	if err != nil {
		t.Errorf("Generated key not a valid hex string")
	}
}

func TestBasicParseConfig(t *testing.T) {
	enMap := make(map[string]string)
	cfg := []byte(`
PORT=8080
URI=localhost
#OLD_PORT=5000
	`)
	err := parseConfig(cfg, enMap, "basic")
	if err != nil {
		t.Errorf("Unable to parse config; %v", err)
	}

	port, ok := enMap["PORT"]
	if !ok {
		t.Errorf("Expected config to have been marshalled into map")
	}

	if port != "8080" {
		t.Errorf("Expected value to be %s, got %s", "8080", port)
	}

	if enMap["OLD_PORT"] != "" {
		t.Errorf("Expected ignored value to not be parsed")
	}

	enMap = make(map[string]string)

	parseConfig(cfg, enMap, "yaml")
	if len(enMap) != 0 {
		t.Errorf("Expected basic style env not be be parseable with yaml envType")
	}
}
func TestYamlParseConfig(t *testing.T) {
	enMap := make(map[string]string)
	cfg := []byte(`
PORT:8080
URI:localhost
#OLD_PORT:5000
	`)
	err := parseConfig(cfg, enMap, "yaml")
	if err != nil {
		t.Errorf("Unable to parse config; %v", err)
	}

	port, ok := enMap["PORT"]
	if !ok {
		t.Errorf("Expected config to have been marshalled into map")
	}

	if port != "8080" {
		t.Errorf("Expected value to be %s, got %s", "8080", port)
	}

	if enMap["OLD_PORT"] != "" {
		t.Errorf("Expected ignored value to not be parsed")
	}
}

func TestCanIgnore(t *testing.T) {
	data := []struct {
		text     string
		expected bool
	}{
		{text: "# url", expected: true},
		{text: "    # url", expected: true},
		{text: "url", expected: false},
		{text: "   url", expected: false},
	}

	for _, val := range data {
		if canIgnore(val.text) != val.expected {
			t.Errorf("Expected canIgnore(%s) to be %v, got %v", val.text, val.expected, canIgnore(val.text))
		}
	}
}

func TestDecodeHex(t *testing.T) {
	_nonce, _text := generateKey(), generateKey()
	hexString := _text + delimiter + _nonce
	_, _, err := decodeFile(hexString)
	if err != nil {
		t.Errorf("Unable to decode valid hex string, got error %v", err)
	}
	nonce, text, err := decodeFile("invalidhex")
	if err == nil {
		t.Errorf("Expected invalid hex file to not decode, got values %s, %s", nonce, text)
	}

}
