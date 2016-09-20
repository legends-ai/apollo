package lib

import (
	"github.com/Sirupsen/logrus"
	"github.com/asunaio/apollo/config"
	"github.com/asunaio/apollo/models"
	"github.com/gocql/gocql"
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

	cluster := gocql.NewCluster(cfg.DBHost...)
	cluster.ProtoVersion = 3
	cluster.Keyspace = cfg.DBKeyspace
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Fatalf("Could not connect to Cassandra: %v", err)
	}
	injector.Map(session)

	// Vulgate
	vulgate, err := models.NewVulgate()
	if err != nil {
		logger.Fatalf("Could not instantiate Vulgate: %v", err)
	}
	_, err = injector.ApplyMap(vulgate)
	if err != nil {
		logger.Fatalf("Could not inject Vulgate: %v", err)
	}

	// Database DAO
	_, err = injector.ApplyMap(models.NewMatchSumDAO())
	if err != nil {
		logger.Fatalf("Could not inject MatchSumDAO: %v", err)
	}

	// Setup aggregator
	_, err = injector.ApplyMap(models.NewDeriver())
	if err != nil {
		logger.Fatalf("Could not inject Deriver: %v", err)
	}

	_, err = injector.ApplyMap(models.NewAggregator())
	if err != nil {
		logger.Fatalf("Could not inject Aggregator: %v", err)
	}

	_, err = injector.ApplyMap(models.NewChampionDAO())
	// _, err = injector.ApplyMap(&models.MockChampionDAO{})
	if err != nil {
		logger.Fatalf("Could not inject ChampionDAO: %v", err)
	}

	return injector
}
