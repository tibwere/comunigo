package main

import (
	"context"
	"log"

	"gitlab.com/tibwere/comunigo/registration/grpchandler"
	"gitlab.com/tibwere/comunigo/utilities"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	err := utilities.InitLogger("registration")
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	cfg, err := utilities.SetupRegistrationServer()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	if cfg.EnableVerbose {
		log.Printf("Start server on port %v (Group size: %v)\n", cfg.RegPort, cfg.ChatGroupSize)
	}

	regServer := grpchandler.NewRegistrationServer(cfg.ChatGroupSize)
	grpcServer := grpc.NewServer()

	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		regServer.UpdateMembers(grpcServer, cfg.SeqHostname, cfg.RegPort, cfg.TypeOfService == "sequencer")
		return nil
	})
	errs.Go(func() error {
		return grpchandler.ServeSignRequests(cfg.RegPort, regServer, grpcServer)
	})

	if err = errs.Wait(); err != nil {
		log.Fatalf("Something went wrong while serving peers (%v)\n", err)
	}

	log.Println("Registration server shutdown")
}
