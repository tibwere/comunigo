package main

import (
	"context"
	"log"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/sequencer/grpchandler"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	membersCh := make(chan string)

	err := config.InitLogger("sequencer")
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	cfg, err := config.SetupSequencer()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	if cfg.EnableVerbose {
		log.Printf("Start server on port %v\n", cfg.ChatPort)
	}

	startupServer := &grpchandler.StartupSequencerServer{
		MembersCh: membersCh,
	}

	seqServer := grpchandler.NewSequencerServer(cfg.ChatPort, cfg.ChatGroupSize)
	grpcServerToGetPeers := grpc.NewServer()

	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		return grpchandler.GetClientsFromRegister(cfg.RegPort, startupServer, grpcServerToGetPeers)
	})
	go errs.Go(func() error {
		seqServer.LoadMembers(membersCh, grpcServerToGetPeers)
		return nil
	})
	if err = errs.Wait(); err != nil {
		log.Fatalf("Something went wrong while retrieving peer list (%v)\n", err)
	}

	errs.Go(func() error {
		seqServer.OrderMessages()
		return nil
	})
	errs.Go(func() error {
		return grpchandler.ServePeers(seqServer)
	})
	if err = errs.Wait(); err != nil {
		log.Fatalf("Something went wrong while sending/receiving messages from peer (%v)\n", err)
	}
}
