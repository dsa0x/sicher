package sicher

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setupTest() (*sicher, string, string) {
	s := New("testenv", "./example")
	return s, fmt.Sprintf("%s%s.enc", s.Path, s.Environment), fmt.Sprintf("%s%s.key", s.Path, s.Environment)

}

func TestNewWithNoEnvironment(t *testing.T) {
	path, _ := filepath.Abs(".")
	path += "/"
	s := New("", "")
	if s.Path != path {
		t.Errorf("Expected path to be %s, got %s", path, s.Path)
	}

	if s.Environment != defaultEnv {
		t.Errorf("Expected environment to be set to %s if none is given, got %s", defaultEnv, s.Environment)
	}

}
func TestNewWithEnvironment(t *testing.T) {
	env := "testenv"
	s := New("testenv", "")

	if s.Environment != env {
		t.Errorf("Expected environment to be set to %s if none is given, got %s", env, s.Environment)
	}

}
func TestEnvStyle(t *testing.T) {
	s := New("testenv", "")
	s.SetEnvStyle("dotenv")

	if s.envStyle != "dotenv" {
		t.Errorf("Expected environment style to be set to %s, got %s", "dotenv", s.envStyle)
	}

}

func TestInvalidEnvStyle(t *testing.T) {
	s := New("testenv", "")
	if os.Getenv("SICHER_ENV_STYLE") == "1" {
		s.SetEnvStyle("wrong")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestInvalidEnvStyle")
	cmd.Env = append(os.Environ(), "SICHER_ENV_STYLE=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, expected exit status 1", err)
}

func TestEditSuccess(t *testing.T) {
	oldExecCmd := execCmd
	defer func() { execCmd = oldExecCmd }()
	s, encPath, keyPath := setupTest()

	s.Initialize(os.Stdin)
	buf := bytes.Buffer{}

	execCmd = func(cmd string, args ...string) *exec.Cmd {
		stdIn, stdOut, stdErr = &buf, &buf, &buf

		if cmd != "vim" {
			t.Errorf("Expected command to be vim, got %s", cmd)
		}
		return exec.Command("cat", args...)
	}

	err := s.Edit("vim")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("TESTKEY=loremipsum")) {
		t.Errorf("Expected credential file to be opened and contain TESTKEY=loremipsum; got %s", buf.String())
	}
	if !bytes.Contains(buf.Bytes(), []byte("File encrypted and saved")) {
		t.Errorf("Expected file to be saved and message to be displayed, got %s", buf.String())
	}

	// get path to the gitignore file and cleanup
	gitPath := strings.Replace(encPath, fmt.Sprintf("%s.enc", s.Environment), ".gitignore", 1)

	t.Cleanup(func() {
		os.Remove(encPath)
		os.Remove(keyPath)
		os.Remove(gitPath)
	})
}
func TestEditFail(t *testing.T) {
	oldExecCmd := execCmd
	defer func() { execCmd = oldExecCmd }()
	s, _, _ := setupTest()

	buf := bytes.Buffer{}

	execCmd = func(cmd string, args ...string) *exec.Cmd {
		stdIn, stdOut, stdErr = &buf, &buf, &buf

		if cmd != "vim" {
			t.Errorf("Expected command to be vim, got %s", cmd)
		}
		return exec.Command("cat", args...)
	}

	err := s.Edit("vim")
	if err == nil {
		t.Errorf("Expected error to be returned, got %s", err)
	}
}

func TestSicherInitialize(t *testing.T) {

	s, encPath, keyPath := setupTest()

	s.Initialize(os.Stdin)

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

}

func TestSicherInitializeExistingCredOverwrite(t *testing.T) {

	s, encPath, keyPath := setupTest()

	f, err := os.Create(encPath)
	if err != nil {
		t.Errorf("Expected credential file to be created; got error %v", err)
	}
	f.Write([]byte("test"))

	buf := bytes.Buffer{}
	buf.WriteString("yes")

	s.Initialize(&buf)

	f, err = os.Open(encPath)
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

	t.Logf("Expects credential file to be overwritten if user confirms with 'yes'")
}

func TestSicherInitializeExistingCredNoOverwrite(t *testing.T) {

	s, encPath, keyPath := setupTest()

	f, err := os.Create(encPath)
	if err != nil {
		t.Errorf("Expected credential file to be created; got error %v", err)
	}
	f.Write([]byte("test"))

	buf := bytes.Buffer{}
	buf.WriteString("n")

	err = s.Initialize(&buf)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	f, err = os.Open(encPath)
	if err != nil {
		t.Errorf("Expected credential file to have been created; got error %v", err)
	}
	f, err = os.Open(keyPath)
	if err == nil {
		t.Errorf("Expected key file to not have been created as user chose not to overwrite")
	}

	// get path to the gitignore file and cleanup
	gitPath := strings.Replace(encPath, fmt.Sprintf("%s.enc", s.Environment), ".gitignore", 1)

	t.Cleanup(func() {
		os.Remove(encPath)
		os.Remove(keyPath)
		os.Remove(gitPath)
		f.Close()
	})

	t.Logf("Expects key file to not have been created as user chose not to overwrite")
}

func TestLoadEnv(t *testing.T) {

	s, encPath, keyPath := setupTest()

	s.Initialize(os.Stdin)

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

// func fakeExecCommand
