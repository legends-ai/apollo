package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// AppConfig represents the application config
type AppConfig struct {
	Port int `required:"true" default:"4834"`
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
