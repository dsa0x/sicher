package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/dsa0x/sicher"
)

var (
	pathFlag          string
	envFlag           string
	editorFlag        string
	styleFlag         string
	gitignorePathFlag string
)

var writer io.Writer = os.Stderr

var errHelp = `
# Initialize sicher in your project
sicher init

# Edit environment variables
sicher edit
`

func init() {
	flag.StringVar(&pathFlag, "path", ".", "Path to the project")
	flag.StringVar(&envFlag, "env", "dev", "Environment to use")
	flag.StringVar(&styleFlag, "style", string(sicher.DefaultEnvStyle), "Env file style. Valid values are dotenv and yaml")
	flag.StringVar(&editorFlag, "editor", "vim", "Select editor.")
	flag.StringVar(&gitignorePathFlag, "gitignore", ".", "Path to the gitignore file")

	flag.ErrHelp = errors.New(errHelp)
	flag.Usage = func() {
		fmt.Fprint(writer, errHelp)
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
	s.SetEnvStyle(styleFlag)
	switch command {
	case "init":
		s.Initialize(os.Stdin)
	case "edit":
		err := s.Edit(editorFlag)
		if err != nil {
			fmt.Fprintln(writer, err)
			os.Exit(1)
		}
	default:
		flag.Usage()
	}
}
