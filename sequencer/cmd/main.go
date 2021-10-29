// Package entry point della logica relativa al sequencer
// all'interno dell'applicazione comuniGO
package main

import (
	"log"
	"os"
	"sync"

	"gitlab.com/tibwere/comunigo/proto"
	"gitlab.com/tibwere/comunigo/sequencer/grpchandler"
	"gitlab.com/tibwere/comunigo/utilities"
	"google.golang.org/grpc"
)

// Launcher del nodo di registrazione
func main() {
	var wg sync.WaitGroup
	ctx := utilities.GetContextForSigHandling()

	err := utilities.InitLogger("sequencer")
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	cfg, err := utilities.InitSequencerConfig()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	// Se la modalità non è sequencer, shutdown!
	if cfg.GetTOS() != utilities.TOS_CS_SEQUENCER {
		log.Printf("Chosen modality do not need sequencer, shutdown!")
		os.Exit(0)
	}

	log.Printf("Start server on port %v\n", cfg.GetToPeersPort())

	// Inizializzazione dei server
	memberCh := make(chan *proto.PeerInfo, cfg.GetMulticastGroupSize())
	seqH := grpchandler.NewSequencerServer(cfg.GetToPeersPort(), cfg.GetMulticastGroupSize(), memberCh)
	fromRegH := grpchandler.NewFromRegisterServer(memberCh)
	fromRegToSeqGRPCserver := grpc.NewServer()

	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := fromRegH.GetPeersFromRegister(ctx, cfg.GetToRegistryPort(), fromRegToSeqGRPCserver); err != nil {
			log.Printf("Unable to retrieve peer list (%v)", err)
		}
	}()

	go func() {
		if err := seqH.StartupConnectionWithPeers(ctx, fromRegToSeqGRPCserver); err != nil {
			log.Printf("Unable to comunicate with peer anymore (%v)", err)
		}
	}()

	go func() {
		if err := seqH.OrderMessages(ctx); err != nil {
			log.Printf("Unable to order message to be delivered (%v)", err)
		}
	}()

	go func() {
		if err := seqH.ServePeers(ctx); err != nil {
			log.Printf("Unable to serve peers anymore (%v)", err)
		}
	}()

	wg.Wait()
}
