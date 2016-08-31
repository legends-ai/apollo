package lib

import (
	"log"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/apollo/config"
	"github.com/simplyianm/apollo/models"
	"github.com/simplyianm/inject"
)

// NewInjector builds a new Injector.
func NewInjector() inject.Injector {
	injector := inject.New()
	injector.Map(injector)

	logger := logrus.New()
	injector.Map(logger)

	cfg := config.Initialize()
	injector.Map(cfg)

	_, err := injector.ApplyMap(&models.ChampionDAO{})
	if err != nil {
		log.Fatalf("Could not inject ChampionDAO: %v", err)
	}

	return injector
}
