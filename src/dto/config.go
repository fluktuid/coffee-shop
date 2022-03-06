package dto

type Config struct {
	BrewDuration int    `envconfig:"brew_duration" default:"1000" required:"true"`
	CustomerWait int    `envconfig:"customer_wait" default:"5000" required:"true"`
	Redis        string `envconfig:"redis" default:"redis" required:"true"`
}
