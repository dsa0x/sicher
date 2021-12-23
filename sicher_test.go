package sicher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTest() (*Sicher, string, string) {
	s := New("testenv", "./example")
	return s, fmt.Sprintf("%s%s.enc", s.Path, s.Environment), fmt.Sprintf("%s%s.key", s.Path, s.Environment)

}

func TestNewWithNoEnvironment(t *testing.T) {
	path, _ := filepath.Abs(".")
	path += "/"
	s := New("")
	if s.Path != path {
		t.Errorf("Expected path to be %s, got %s", path, s.Path)
	}

	if s.Environment != defaultEnv {
		t.Errorf("Expected environment to be set to %s if none is given, got %s", defaultEnv, s.Environment)
	}

}
func TestNewWithEnvironment(t *testing.T) {
	env := "testenv"
	s := New("testenv")

	if s.Environment != env {
		t.Errorf("Expected environment to be set to %s if none is given, got %s", env, s.Environment)
	}

}

func TestSicherInitialize(t *testing.T) {

	s, encPath, keyPath := setupTest()

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

	s, encPath, keyPath := setupTest()

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

func TestSetEnv(t *testing.T) {
	s, _, _ := setupTest()

	s.data["PORT"] = "8080"
	s.setEnv()

	if os.Getenv("PORT") != "8080" {
		t.Errorf("Expected environment variable %s to have been set to %s, got %s", "PORT", "8080", os.Getenv("PORT"))
	}

}
