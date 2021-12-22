package example

import (
	"fmt"

	"github.com/dsa0x/sicher/sicher"
)

type Config struct {
	Port        string `required:"false" envconfig:"PORT"`
	MongoDbURI  string `required:"true" envconfig:"MONGO_DB_URI"`
	MongoDbName string `required:"true" envconfig:"MONGO_DB_NAME"`
	JWTSecret   string `required:"false" envconfig:"JWT_SECRET"`
}

func Configure() {

	var cfg Config
	// cfg := make(map[string]string)

	s := sicher.New("dev")
	err := s.LoadEnv("", &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg)
}
