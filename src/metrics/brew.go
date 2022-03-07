package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	brew = promauto.NewCounter(prometheus.CounterOpts{
		Name: "coffeeshop_brews_total",
		Help: "The total number of brewed coffees",
	})
	brew_sort = make(map[string]prometheus.Counter)
	order     = promauto.NewCounter(prometheus.CounterOpts{
		Name: "coffeeshop_order_total",
		Help: "The total number of brewed coffees",
	})
	order_sort = make(map[string]prometheus.Counter)
	successful = promauto.NewCounter(prometheus.CounterOpts{
		Name: "coffeeshop_successful",
		Help: "The total number of brewed coffees",
	})
	unsuccessful = promauto.NewCounter(prometheus.CounterOpts{
		Name: "coffeeshop_unsuccessful",
		Help: "The total number of brewed coffees",
	})
	spawn = promauto.NewCounter(prometheus.CounterOpts{
		Name: "coffeeshop_spawn_total",
		Help: "The total number of brewed coffees",
	})
	spawnGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "coffeeshop_spawn",
		Help: "The total number of brewed coffees",
	})
)

func Brew(sort string) {
	brew.Inc()
	sort_ctr := brew_sort[sort]
	if sort_ctr == nil {
		sort_ctr := promauto.NewCounter(prometheus.CounterOpts{
			Name: "coffeeshop_brews_" + sort,
			Help: "The total number of brewed " + sort,
		})
		brew_sort[sort] = sort_ctr
	}
	sort_ctr.Inc()
}

func Order(sort string) {
	order.Inc()
	sort_ctr := order_sort[sort]
	if sort_ctr == nil {
		sort_ctr := promauto.NewCounter(prometheus.CounterOpts{
			Name: "coffeeshop_order_" + sort,
			Help: "The total number of brewed " + sort,
		})
		order_sort[sort] = sort_ctr
	}
	sort_ctr.Inc()
}

func Customer(success bool) {
	if success {
		successful.Inc()
	} else {
		unsuccessful.Inc()
	}
}

func Spawn(cnt int) {
	spawn.Add(float64(cnt))
	spawnGauge.Set(float64(cnt))
}
