package util

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

func LoadEnv(i interface{}) {
	err := envconfig.Process("", i)
	if err != nil {
		log.Fatal(err)
	}
}
