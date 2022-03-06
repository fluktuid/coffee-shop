package main

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fluktuid/coffee-shop/src/dto"
	"github.com/fluktuid/coffee-shop/src/util"
)

func main() {
	var config dto.Config

	util.LoadEnv(&config)

	redis := util.NewClient(config.Redis)
	go func() {
		for {
			orderStr, _ := redis.AnyPattern(util.ORDER_PREFIX + "*")
			if orderStr == "" {
				time.Sleep(time.Duration(50) + time.Millisecond)
				continue
			}
			logrus.Infof("Got Order")
			var order dto.Coffee
			util.Unmarshal(orderStr, &order)
			if order.Sort != "" && order.Customer != "" {
				logrus.Infof(" - trying to brew %s\n", order.Sort)
				redis.LPush(util.BREW_IN, order.Sort)
				coffee, err := redis.BRPop(util.BREW_OUT + order.Sort)
				if coffee != "" && err == nil {
					redis.Set(util.OBTAIN_PREFIX+order.Customer, coffee, time.Duration(config.CustomerWait)*time.Millisecond)
				}
			}
		}
	}()

	select {}
}
