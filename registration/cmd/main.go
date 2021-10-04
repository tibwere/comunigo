package main

import (
	"log"
	"sync"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/registration/grpchandler"
	"google.golang.org/grpc"
)

func main() {
	var wg sync.WaitGroup

	config, err := config.SetupRegistrationServer()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	if config.EnableVerbose {
		log.Printf("Start server on port %v (Group size: %v)\n", config.RegPort, config.ChatGroupSize)
	}

	regServer := grpchandler.NewRegistrationServer(config.ChatGroupSize)
	grpcServer := grpc.NewServer()

	wg.Add(2)
	go regServer.UpdateMembers(grpcServer, config.SeqHostname, config.RegPort, &wg)
	go grpchandler.ServeSignRequests(config.RegPort, regServer, grpcServer, &wg)
	wg.Wait()
	log.Println("Registration server shutdown")
}
