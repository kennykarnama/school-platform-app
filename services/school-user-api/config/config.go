package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RestPort                string `envconfig:"PORT" default:"8080"`
	ServiceName             string `envconfig:"SERVICE_NAME" default:"user_api"`
	AccessTokenCfg          AccessTokenConfig
	RefreshTokenCfg         RefreshTokenConfig
	EnableEmailConfirmation bool `envconfig:"ENABLE_EMAIL_CONFIRMATION" default:"false"`
}

type AccessTokenConfig struct {
	Secret     string        `envconfig:"ACCESS_TOKEN_SECRET" required:"true"`
	Expiration time.Duration `envconfig:"ACCESS_TOKEN_EXPIRATION" default:"15m"`
}

type RefreshTokenConfig struct {
	Secret     string        `envconfig:"REFRESH_TOKEN_SECRET" required:"true"`
	Expiration time.Duration `envconfig:"REFRESH_TOKEN_EXPIRATION" default:"60m"`
}

type MailGunConfig struct {
	UseMailGun    bool   `envconfig:"USE_MAIL_GUN"`
	EmailDomain   string `envconfig:"EMAIL_DOMAIN"`
	MailGunApiKey string `envconfig:"MAIL_GUN_API_KEY"`
}

func Get() Config {
	var cfg Config
	envconfig.MustProcess("", &cfg)
	return cfg
}
