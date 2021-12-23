package cli

import (
	"bytes"
	"os"
	"testing"
)

func TestInvalidCmd(t *testing.T) {
	oldWriter := writer
	defer func() { writer = oldWriter }()
	b := bytes.Buffer{}

	writer = &b
	os.Args = []string{"sicher", "", "testenv"}
	Execute()
	if !bytes.Equal(b.Bytes(), []byte(errHelp)) {
		t.Fatalf("Expected to print help text %s if command is invalid, got %s", errHelp, b.String())
	}
}

func TestInitCmd(t *testing.T) {
	oldWriter := writer
	defer func() { writer = oldWriter }()
	b := bytes.Buffer{}

	writer = &b
	os.Args = []string{"sicher", "init", "testenv"}
	Execute()
	if bytes.Equal(b.Bytes(), []byte(errHelp)) {
		t.Fatalf("Expected to not print help text if command is valid")
	}

	t.Cleanup(func() {
		os.Remove("dev.enc")
		os.Remove("dev.key")
		os.Remove(".gitignore")
	})
}
