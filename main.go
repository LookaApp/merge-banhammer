package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"mergeban/pkg"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	notifier := mergeban.NewSlackNotifier()

	listenAddress := os.Getenv("MERGEBAN_LISTEN_ADDR")
	if listenAddress == "" {
		logger.Fatalf("Was not supplied with required environment variable MERGEBAN_LISTEN_ADDR\n")
	}

	banService := mergeban.CreateService(logger, notifier)

	httpServer := http.Server{
		ErrorLog: logger,
		Handler:  banService,
	}

	listener, err := net.Listen("tcp4", listenAddress)
	if err != nil {
		logger.Fatalf("Failed to listen on %s: %v\n", listenAddress, err)
	}

	err = httpServer.Serve(listener)
	if err != nil { // Always non-nil - blocks until server shutsdown
		logger.Fatalf("Stopping server: %v\n", err)
	}
}
