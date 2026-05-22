package core_auth

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Secret string        `envconfig:"SECRET" required:"true"`
	Expiry time.Duration `envconfig:"EXPIRY" default:"24h"`
	Issuer string        `envconfig:"ISSUER" default:"todo-web"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("JWT", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	if config.Secret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required: %w", fmt.Errorf("process envconfig"))
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("new JWT config: %w", err))
	}

	return config
}
