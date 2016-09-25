package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/simplyianm/inject"
	"google.golang.org/grpc"

	"github.com/asunaio/apollo/config"
	"github.com/asunaio/apollo/lib"
	"github.com/asunaio/apollo/server"

	_ "net/http/pprof"

	apb "github.com/asunaio/apollo/gen-go/asuna"
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

	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		monitorPort := fmt.Sprintf(":%d", config.MonitorPort)
		logger.Infof("Monitor listening on %s", monitorPort)
		http.ListenAndServe(monitorPort, nil)
	}()

	apb.RegisterApolloServer(s, serv)
	logger.Infof("Listening on %s", port)
	s.Serve(lis)
}
