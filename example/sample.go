package main

import (
	"fmt"

	"github.com/dsa0x/sicher"
)

type Config struct {
	Port        string `required:"false" env:"PORT"`
	MongoDbURI  string `required:"false" env:"MONGO_DB_URI"`
	MongoDbName string `required:"false" env:"MONGO_DB_NAME"`
	TestKey     string `required:"false" env:"TESTKEY"`
}

// LoadConfigStruct Loads config into a struct
func LoadConfigStruct() {

	var cfg Config

	s := sicher.New("dev", ".")
	err := s.LoadEnv("", &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg)
}

// LoadConfigMap Loads config into a map
func LoadConfigMap() {

	cfg := make(map[string]string)

	s := sicher.New("dev", ".")
	s.SetEnvStyle("yaml") // default is dotenv
	err := s.LoadEnv("", &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg)
}

func main() {
	LoadConfigStruct()
	LoadConfigMap()
}
