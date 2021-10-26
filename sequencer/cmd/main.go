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

func main() {
	var wg sync.WaitGroup
	ctx := utilities.GetContextForSigHandling()

	// Inizializzazione dell'attività di logging su file dedicato
	err := utilities.InitLogger("sequencer")
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	// Retrieve delle impostazioni di configurazione dall'ambiente
	cfg, err := utilities.SetupSequencer()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	// Se la modalità non è sequencer, shutdown!
	if cfg.TypeOfService != "sequencer" {
		log.Printf("Chosen modality do not need sequencer, shutdown!")
		os.Exit(0)
	}

	log.Printf("Start server on port %v\n", cfg.ChatPort)

	// Inizializzazione dei server
	memberCh := make(chan *proto.PeerInfo, cfg.ChatGroupSize)
	seqH := grpchandler.NewSequencerServer(cfg.ChatPort, cfg.ChatGroupSize, memberCh)
	fromRegH := grpchandler.NewFromRegisterServer(memberCh)
	fromRegToSeqGRPCserver := grpc.NewServer()

	wg.Add(4)

	// Retrieve dei peer dal register
	go func() {
		defer wg.Done()
		if err := fromRegH.GetPeersFromRegister(ctx, cfg.RegPort, fromRegToSeqGRPCserver); err != nil {
			log.Printf("Unable to retrieve peer list (%v)", err)
		}
	}()
	// Metodo buffer che non fa altro che prendere da un canale
	// degli indirizzi ed utilizzarli per aprire nuove connessioni
	go func() {
		if err := seqH.StartupConnectionWithPeers(ctx, fromRegToSeqGRPCserver); err != nil {
			log.Printf("Unable to comunicate with peer anymore (%v)", err)
		}
	}()

	// Routine per l'ordinamento sequenziale dei messaggi
	go func() {
		if err := seqH.OrderMessages(ctx); err != nil {
			log.Printf("Unable to order message to be delivered (%v)", err)
		}
	}()

	// GRPC server per servire i peers con messaggi provenienti
	// dalla routine precedente
	go func() {
		if err := seqH.ServePeers(ctx); err != nil {
			log.Printf("Unable to serve peers anymore (%v)", err)
		}
	}()

	wg.Wait()
}
