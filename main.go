package main

import (
	"fmt"
	"os"

	sicher "github.com/dsa0x/sicher/cmd"
)

func main() {
	sicher.Execute()
	cfg := sicher.Config{}
	cfg.Get()
	cfg.SetEnv()

	TestEnv()
}

func TestEnv() {
	fmt.Println(os.Getenv("mongodb"), "ENV")
}
