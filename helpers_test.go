package sicher

import (
	"encoding/hex"
	"fmt"
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
	err := parseConfig(cfg, enMap, "dotenv")
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
		t.Errorf("Expected dotenv style env not be be parseable with yaml envType")
	}
}
func TestParseConfig(t *testing.T) {

	tests := []struct {
		text     string
		expected map[string]string
		envType  string
	}{
		{
			text: `
			PORT:8080
			URI:localhost
			#OLD_PORT:5000
				`,
			expected: map[string]string{
				"PORT": "8080",
				"URI":  "localhost",
			},
			envType: "yaml",
		},
		{
			text: `
			PORT=8080
			URI=localhost
			#OLD_PORT=5000
				`,
			expected: map[string]string{
				"PORT": "8080",
				"URI":  "localhost",
			},
			envType: "dotenv",
		},
		{
			text: `
			PORT=8080
			URI=localhost
			#OLD_PORT=5000
			KEY=value=ndsjhjdghdhg
				`,
			expected: map[string]string{
				"PORT": "8080",
				"URI":  "localhost",
				"KEY":  "value=ndsjhjdghdhg",
			},
			envType: "dotenv",
		},
		{
			text: `
			PORT:8080
			URI=localhost
				`,
			expected: map[string]string{
				"PORT": "8080",
			},
			envType: "yaml",
		},
		{
			text: `
			PORT:8080
			URI=localhost
			SOME_KEY:somevalue=jsfhjdghdhg
				`,
			expected: map[string]string{
				"URI": "localhost",
			},
			envType: "dotenv",
		},
	}

	for _, val := range tests {
		enMap := make(map[string]string)
		if err := parseConfig([]byte(val.text), enMap, EnvStyle(val.envType)); err != nil {
			t.Errorf("Unable to parse config; %v", err)
		}

		t.Run(fmt.Sprintf("Envtype %s", val.envType), func(t *testing.T) {
			for key, value := range val.expected {
				if enMap[key] != value {
					t.Errorf("Expected value to be %s, got %s", value, enMap[key])
				}
			}
			for key, value := range enMap {
				if val.expected[key] != value {
					t.Errorf("Expected value to be %s, got %s", value, val.expected[key])
				}
			}
		})
	}

}

func TestYamlParseConfigError(t *testing.T) {
	enMap := make(map[string]string)
	cfg := []byte(`
PORT:8080
URI:localhost
#OLD_PORT:5000
	`)
	err := parseConfig(cfg, enMap, "wrong")
	if err == nil {
		t.Errorf("Expected error to be thrown when parsing wrong envType")
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
