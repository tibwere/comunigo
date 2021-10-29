// Package entry point della logica relativa al nodo di
// registrazione all'interno dell'applicazione comuniGO
package main

import (
	"log"
	"sync"

	"gitlab.com/tibwere/comunigo/registration/grpchandler"
	"gitlab.com/tibwere/comunigo/utilities"
	"google.golang.org/grpc"
)

// Launcher del nodo di registrazione
func main() {
	var wg sync.WaitGroup
	ctx := utilities.GetContextForSigHandling()

	err := utilities.InitLogger("registration")
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	cfg, err := utilities.InitRegistrationServiceConfig()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	log.Printf(
		"Start server on port %v (Group size: %v)\n",
		cfg.GetExposedPort(),
		cfg.GetMulticastGroupSize(),
	)

	regServer := grpchandler.NewRegistrationServer(cfg.GetMulticastGroupSize())
	grpcServer := grpc.NewServer()

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := regServer.UpdateMembers(
			ctx,
			grpcServer,
			cfg.GetSequencerAddress(),
			cfg.GetExposedPort(),
			cfg.GetTOS() == utilities.TOS_CS_SEQUENCER,
		)
		if err != nil {
			log.Printf("Unable to update members (%v)", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := grpchandler.ServeSignRequests(ctx, cfg.GetExposedPort(), regServer, grpcServer)
		if err != nil {
			log.Printf("Unable to serve sign requests (%v)", err)
		}
	}()

	wg.Wait()

	log.Println("Registration server shutdown")
}
