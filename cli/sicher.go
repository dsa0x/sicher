package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/dsa0x/sicher/sicher"
)

var (
	pathFlag   string
	envFlag    string
	editorFlag string
)

func init() {
	errHelp := `
# Initialize sicher in your project
sicher init

# Edit environment variables
sicher edit
	`
	flag.StringVar(&pathFlag, "path", ".", "Path to the project")
	flag.StringVar(&envFlag, "env", "dev", "Environment to use")
	flag.StringVar(&editorFlag, "editor", "vim", "Select editor. vim | vi | nano")

	flag.ErrHelp = errors.New(errHelp)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, errHelp)
	}
}

func Execute() {
	flag.Parse()
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	command := os.Args[1]
	flag.CommandLine.Parse(os.Args[2:])
	s := sicher.New(envFlag, pathFlag)
	switch command {
	case "init":
		s.Initialize()
	case "edit":
		s.Edit(editorFlag)
	}
}
