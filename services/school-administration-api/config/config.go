package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RestPort            string `envconfig:"PORT" default:"8080"`
	ServiceName         string `envconfig:"SERVICE_NAME" default:"school-administration-api"`
	SessionTTL          int    `envconfig:"SESSION_TTL" default:"300"`
	EnableAuth          bool   `envconfig:"ENABLE_AUTH" default:"false"`
	SessionCookieSecure bool   `envconfig:"SESSION_COOKIE_SECURE" default:"false"`
}

func Get() Config {
	var cfg Config
	envconfig.MustProcess("", &cfg)
	return cfg
}
