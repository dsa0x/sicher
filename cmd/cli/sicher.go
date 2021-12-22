package main

import (
	"github.com/dsa0x/sicher/cli"
	"github.com/dsa0x/sicher/sicher"
)

// func main() {
// 	cli.Execute()
// }

// use default values
var sich = &sicher.Sicher{Environment: "dev", Path: "."}
var (
	pathFlag   string
	envFlag    string
	editorFlag string
)

func main() {
	cli.Execute()
}
