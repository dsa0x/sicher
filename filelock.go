package sicher

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

type fileLock struct {
	file *os.File
}

type FileCreator func() (*os.File, error)

func newFileLock(file *os.File) *fileLock {
	return &fileLock{file: file}
}

func (l *fileLock) Lock() bool {
	if err := unix.Flock(int(l.file.Fd()), unix.LOCK_EX); err != nil {
		log.Fatalf("file lock error: %v\n", err)
		return false
	}
	return true
}

func (l *fileLock) Unlock() {
	if err := unix.Flock(int(l.file.Fd()), unix.LOCK_UN); err != nil {
		log.Fatalf("file unlock error: %v\n", err)
	}
}

func (l *fileLock) Create(creator FileCreator) error {
	f, err := creator()
	if err != nil {
		return err
	}
	defer f.Close()

	l.file = f
	l.Lock()
	return nil
}

func (l *fileLock) LockWithTimeout(timeout time.Duration) {
	ch := make(chan bool)
	go func() {
		ch <- l.Lock()
	}()

	select {
	case <-time.After(timeout):
		fmt.Printf("File is being edited in another terminal. Lock timeout after %v\n", timeout)
		os.Exit(1)
	case <-ch:
		break
	}
}
