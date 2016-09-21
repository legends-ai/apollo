package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// AppConfig represents the application config
type AppConfig struct {
	Port       int      `required:"true" default:"4834"`
	HealthPort int      `required:"true" default:"9834"`
	DBHost     []string `default:"127.0.0.1"`
	DBKeyspace string   `default:"athena_out"`
}

// Initialize initializes the configuration from env vars
func Initialize() *AppConfig {
	var cfg AppConfig
	err := envconfig.Process("apollo", &cfg)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	return &cfg
}
