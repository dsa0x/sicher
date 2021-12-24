package sicher

import (
	"os"
	"testing"
)

func TestNewFileLock(t *testing.T) {
	f, err := os.CreateTemp("", "test")
	if err != nil {
		t.Errorf("Error when opening test file: %v", err)
		return
	}

	fl := newFileLock(f)
	if fl.file != f {
		t.Errorf("Expected file lock to have file %v; got %v", f, fl.file)
	}

	t.Cleanup(func() {
		f.Close()
		os.Remove(f.Name())
	})
}
