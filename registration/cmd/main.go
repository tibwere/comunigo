package main

import (
	"context"
	"log"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/registration/grpchandler"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	err := config.InitLogger("registration")
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	cfg, err := config.SetupRegistrationServer()
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
		regServer.UpdateMembers(grpcServer, cfg.SeqHostname, cfg.RegPort)
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
