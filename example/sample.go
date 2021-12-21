package example

import (
	"fmt"
	"os"

	"github.com/dsa0x/sicher/sicher"
)

func GetEnv() {
	// cli.Execute()
	// cli.LoadEnv()
	env := "development"
	path := "."
	s := sicher.New(env, path)
	s.LoadEnv()

	fmt.Println(os.Getenv("mongodb"), "ENV")
}
