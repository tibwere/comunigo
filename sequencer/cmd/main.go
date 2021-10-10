package main

import (
	"context"
	"log"
	"os"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/sequencer/grpchandler"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	membersCh := make(chan string)

	// Inizializzazione dell'attività di logging su file dedicato
	err := config.InitLogger("sequencer")
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	// Retrieve delle impostazioni di configurazione dall'ambiente
	cfg, err := config.SetupSequencer()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	// Se la modalità non è sequencer, shutdown!
	if cfg.TypeOfService != "sequencer" {
		log.Printf("Chosen modality do not need sequencer, shutdown!")
		os.Exit(0)
	}

	if cfg.EnableVerbose {
		log.Printf("Start server on port %v\n", cfg.ChatPort)
	}

	// TODO MIGLIORAREEEE
	startupServer := &grpchandler.StartupSequencerServer{
		MembersCh: membersCh,
	}

	seqServer := grpchandler.NewSequencerServer(cfg.ChatPort, cfg.ChatGroupSize)
	grpcServerToGetPeers := grpc.NewServer()

	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		return grpchandler.GetClientsFromRegister(cfg.RegPort, startupServer, grpcServerToGetPeers)
	})
	errs.Go(func() error {
		return seqServer.LoadMembers(membersCh, grpcServerToGetPeers)
	})
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
