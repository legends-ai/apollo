package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/inject"
	"google.golang.org/grpc"

	"github.com/simplyianm/apollo/config"
	"github.com/simplyianm/apollo/lib"
	"github.com/simplyianm/apollo/server"

	apb "github.com/simplyianm/apollo/gen-go/apollo"
)

func main() {
	injector := lib.NewInjector()
	_, err := injector.Invoke(initServer)
	if err != nil {
		log.Fatalf("Could not invoke init: %v", err)
	}
}

func initServer(injector inject.Injector, logger *logrus.Logger, config *config.AppConfig) {
	// Listen on port
	port := fmt.Sprintf(":%d", config.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Setup gRPC server
	s := grpc.NewServer()
	serv := &server.Server{}

	_, err = injector.ApplyMap(serv)
	if err != nil {
		logger.Fatalf("Could not inject server: %v", err)
	}

	apb.RegisterApolloServer(s, serv)
	logger.Infof("Listening on %s", port)
	s.Serve(lis)
}
