package dto

type Config struct {
	BrewDuration int    `envconfig:"brew_duration" default:"1000" required:"true"`
	CustomerWait int    `envconfig:"customer_wait" default:"5000" required:"true"`
	Redis        string `envconfig:"redis" default:"localhost:6379" required:"true"`
	SpawnMax     int    `envconfig:"spawn_max" default:"4" required:"true"`
	CustomerImg  string `envconfig:"customer_img" default:"ghcr.io/fluktuid/coffee-shop/customer" required:"true"`
}
