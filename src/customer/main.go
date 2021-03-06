package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fluktuid/coffee-shop/src/dto"
	"github.com/fluktuid/coffee-shop/src/metrics"
	"github.com/fluktuid/coffee-shop/src/util"
)

func main() {
	var config dto.Config
	util.LoadEnv(&config)
	go metrics.StartMetricsBlock()

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	redis := util.NewClient(config.Redis)
	key := util.ORDER_PREFIX + strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + hostname
	order := dto.Coffee{
		Sort:     "coffee",
		Customer: hostname,
	}
	wait := time.Duration(config.CustomerWait) * time.Millisecond
	redis.Set(key, util.Marshal(order), 0)

	time.Sleep(wait)

	resp, err := redis.GetDel(util.OBTAIN_PREFIX + hostname)

	if err != nil || resp == "" {
		log.Warn("Order not successful")
		metrics.Customer(false)

		// sleep for metrics scraper
		time.Sleep(time.Duration(2) * time.Second)
		log.Exit(1)
	} else {
		metrics.Customer(true)
		log.Infof("Got %s\n", resp)

		// sleep for metrics scraper
		time.Sleep(time.Duration(2) * time.Second)
	}
}
