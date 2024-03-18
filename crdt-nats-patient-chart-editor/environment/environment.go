package environment

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	NATS_URL string `default:"my-nats:4222"`
}

func GetEnvironment() *Environment {
	var env Environment
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &env
}
