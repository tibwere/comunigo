package main

import (
	"context"
	"log"
	"os"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/sequencer/grpchandler"
	"golang.org/x/sync/errgroup"
)

func main() {

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

	// Inizializzazione del sequence server
	seqH := grpchandler.NewSequencerServer(cfg.ChatPort, cfg.ChatGroupSize)

	// Retrieve dei peer dal register
	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		return seqH.GetPeersFromRegister(cfg.RegPort)
	})

	// Routine per l'ordinamento sequenziale dei messaggi
	errs.Go(func() error {
		seqH.OrderMessages()
		return nil
	})

	// GRPC server per servire i peers con messaggi provenienti
	// dalla routine precedente
	errs.Go(func() error {
		return seqH.ServePeers()
	})
	if err = errs.Wait(); err != nil {
		log.Fatalf("Something went wrong while sending/receiving messages from peer (%v)\n", err)
	}
}
