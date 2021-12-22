package sicher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTest(path string) (*Sicher, string, string) {

	if path == "" {
		path = "."
	}
	path, _ = filepath.Abs(path)
	s := &Sicher{
		Path:        path,
		Environment: "testenv",
	}

	return s, fmt.Sprintf("%s/%s.enc", path, s.Environment), fmt.Sprintf("%s/%s.key", path, s.Environment)

}

func TestSicherInitialization(t *testing.T) {

	s, encPath, keyPath := setupTest("../example")

	s.Initialize()

	f, err := os.Open(encPath)
	if err != nil {
		t.Errorf("Expected credential file to have been created; got error %v", err)
	}
	f, err = os.Open(keyPath)
	if err != nil {
		t.Errorf("Expected key file to have been created; got error %v", err)
	}

	// get path to the gitignore file and cleanup
	gitPath := strings.Replace(encPath, fmt.Sprintf("%s.enc", s.Environment), ".gitignore", 1)

	t.Cleanup(func() {
		os.Remove(encPath)
		os.Remove(keyPath)
		os.Remove(gitPath)
		f.Close()
	})

	t.Logf("Credential and key file have been created successfully")

}

func TestLoadEnv(t *testing.T) {

	s, encPath, keyPath := setupTest("../example")

	s.Initialize()

	f, err := os.Open(encPath)
	if err != nil {
		t.Errorf("Expected credential file to have been created; got error %v", err)
	}
	f, err = os.Open(keyPath)
	if err != nil {
		t.Errorf("Expected key file to have been created; got error %v", err)
	}

	mp := make(map[string]string)
	err = s.LoadEnv("", &mp)
	if err != nil {
		t.Errorf("Expected to load envirnoment variables; got error %v", err)
	}

	if len(mp) != 1 {
		t.Errorf("Expected config file to be been populated with env variables")
	}

	// get path to the gitignore file and cleanup
	gitPath := strings.Replace(encPath, fmt.Sprintf("%s.enc", s.Environment), ".gitignore", 1)

	t.Cleanup(func() {
		os.Remove(encPath)
		os.Remove(keyPath)
		os.Remove(gitPath)
		f.Close()
	})

}
