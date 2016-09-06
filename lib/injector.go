package lib

import (
	"log"

	"github.com/Sirupsen/logrus"
	"github.com/gocql/gocql"
	"github.com/simplyianm/apollo/aggregation"
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

	cluster := gocql.NewCluster(cfg.DBHost...)
	cluster.Keyspace = cfg.DBKeyspace
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Could not connect to Cassandra: %v", err)
	}
	injector.Map(session)

	// Setup aggregator
	injector.ApplyMap(&aggregation.Aggregator{})

	return injector
}
