package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fluktuid/coffee-shop/src/dto"
	"github.com/fluktuid/coffee-shop/src/util"
)

func main() {
	var config dto.Config
	util.LoadEnv(config)

	redis := util.NewClient(config.Redis)
	go func() {
		for {
			val, _ := redis.BRPop(util.BREW_IN)
			if val == "" {
				continue
			}
			log.Infof("brewing %s (dur: %d millis)", val, config.BrewDuration)
			time.Sleep(time.Duration(config.BrewDuration))
			log.Infof(" -> ready\n", config.BrewDuration)
			redis.LPush(util.BREW_OUT+val, val)
		}
	}()

	select {}
}