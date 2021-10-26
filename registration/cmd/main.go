package main

import (
	"log"
	"sync"

	"gitlab.com/tibwere/comunigo/registration/grpchandler"
	"gitlab.com/tibwere/comunigo/utilities"
	"google.golang.org/grpc"
)

func main() {
	var wg sync.WaitGroup
	ctx := utilities.GetContextForSigHandling()

	err := utilities.InitLogger("registration")
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	cfg, err := utilities.SetupRegistrationServer()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	log.Printf("Start server on port %v (Group size: %v)\n", cfg.RegPort, cfg.ChatGroupSize)

	regServer := grpchandler.NewRegistrationServer(cfg.ChatGroupSize)
	grpcServer := grpc.NewServer()

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := regServer.UpdateMembers(ctx, grpcServer, cfg.SeqHostname, cfg.RegPort, cfg.TypeOfService == "sequencer"); err != nil {
			log.Printf("Unable to update members (%v)", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := grpchandler.ServeSignRequests(ctx, cfg.RegPort, regServer, grpcServer); err != nil {
			log.Printf("Unable to serve sign requests (%v)", err)
		}
	}()

	wg.Wait()

	log.Println("Registration server shutdown")
}
