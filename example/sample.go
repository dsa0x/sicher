package example

import (
	"fmt"

	"github.com/dsa0x/sicher/sicher"
)

type Config struct {
	Port        string `required:"true" envconfig:"PORT"`
	MongoDbURI  string `required:"true" envconfig:"MONGO_DB_URI"`
	MongoDbName string `required:"true" envconfig:"MONGO_DB_NAME"`
	JWTSecret   string `required:"false" envconfig:"JWT_SECRET"`
}

func Configure() {
	var config Config

	s := sicher.New("development")
	err := s.LoadEnv("REG", &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config.Port)
}
