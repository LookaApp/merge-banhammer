package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"mergeban/mergeban"
)

func main() {
	listenAddress := "localhost:1337"
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	banService := mergeban.CreateService(logger)

	httpServer := http.Server{
		ErrorLog: logger,
		Handler:  banService,
	}

	listener, err := net.Listen("tcp4", listenAddress)
	if err != nil { // Always non-nil - blocks until server shutsdown
		logger.Fatalf("Failed to listen on %s: %v\n", listenAddress, err)
	}

	err = httpServer.Serve(listener)
	if err != nil { // Always non-nil - blocks until server shutsdown
		logger.Fatalf("Stopping server: %v\n", err)
	}
}
