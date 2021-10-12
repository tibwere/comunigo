package main

import (
	"context"
	"log"
	"os"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/proto"
	"gitlab.com/tibwere/comunigo/sequencer/grpchandler"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
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

	// Inizializzazione dei server
	memberCh := make(chan *proto.PeerInfo, cfg.ChatGroupSize)
	seqH := grpchandler.NewSequencerServer(cfg.ChatPort, cfg.ChatGroupSize, memberCh)
	fromRegH := grpchandler.NewFromRegisterServer(memberCh)
	fromRegToSeqGRPCserver := grpc.NewServer()

	// Retrieve dei peer dal register
	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		return fromRegH.GetPeersFromRegister(cfg.RegPort, fromRegToSeqGRPCserver)
	})

	// Metodo buffer che non fa altro che prendere da un canale
	// degli indirizzi ed utilizzarli per aprire nuove connessioni
	errs.Go(func() error {
		return seqH.StartupConnectionWithPeers(fromRegToSeqGRPCserver)
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
