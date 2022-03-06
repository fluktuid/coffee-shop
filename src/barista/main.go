package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fluktuid/coffee-shop/src/dto"
	"github.com/fluktuid/coffee-shop/src/util"
)

func main() {
	var config dto.Config

	redis := util.NewClient(config.Redis)
	go func() {
		for {
			orderStr, _ := redis.AnyPattern(util.ORDER_PREFIX + "*")
			if orderStr == "" {
				time.Sleep(time.Duration(50))
				continue
			}
			var order dto.Coffee
			util.Unmarshal(orderStr, order)
			if order.Sort != "" {
				redis.LPush(util.BREW_IN, order.Sort)
				coffee, err := redis.BRPop(util.BREW_OUT + order.Sort)
				if coffee != "" && err == nil {
					redis.Set(util.OBTAIN_PREFIX+order.Customer, coffee, time.Duration(config.CustomerWait))
				}
			}
		}
	}()

	select {}
	util.LoadEnv(config)

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		log.Info("ordered some coffee")
		fmt.Fprintf(w, "Hello!\n")
		keys, ok := r.URL.Query()["type"]
		if !ok {
			keys = []string{"coffee"}
		}

		// 'brew' coffee
		for _, e := range keys {
			brew(config, e)
			fmt.Fprintf(w, "Here, your %s\n", e)
		}
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func brew(c dto.Config, _type string) string {
	c.Machine
}
