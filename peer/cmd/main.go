package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/peer/grpchandler/p2p/scalar"
	"gitlab.com/tibwere/comunigo/peer/grpchandler/p2p/vectorial"
	"gitlab.com/tibwere/comunigo/peer/grpchandler/reg"
	"gitlab.com/tibwere/comunigo/peer/grpchandler/seq"
	"gitlab.com/tibwere/comunigo/peer/webserver"
	"gitlab.com/tibwere/comunigo/utilities"
	"golang.org/x/sync/errgroup"
)

func main() {
	var wg sync.WaitGroup
	ctx := utilities.GetContextForSigHandling()

	rand.Seed(time.Now().UnixNano())

	cfg, err := utilities.SetupPeer()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	status, err := peer.Init(cfg.RedisHostname)
	if err != nil {
		log.Fatalf("Unable to initialize status (%v)\n", err)
	}

	err = utilities.InitLogger(fmt.Sprintf("peer_%v_main", status.PublicIP))
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	ws := webserver.New(cfg.WebServerPort, cfg.ChatGroupSize, cfg.TypeOfService, status)

	wg.Add(2)
	go ws.Startup(ctx, &wg)
	go internalLogic(ctx, cfg, status, &wg)
	wg.Wait()

	log.Println("Peer is shutting down")
}

func internalLogic(ctx context.Context, cfg *utilities.PeerConfig, status *peer.Status, wg *sync.WaitGroup) {
	defer wg.Done()

	regH := reg.NewToRegisterGRPCHandler(cfg.RegHostname, cfg.RegPort, cfg.EnableVerbose, status)
	err := regH.SignToRegister(ctx)
	if err != nil {
		log.Printf("Unable to sign to register node (%v)\n", err)
		return
	}

	switch cfg.TypeOfService {
	case "sequencer":
		sequencerHandler(ctx, cfg.SeqHostname, cfg.ChatPort, cfg.EnableVerbose, status)
	case "scalar":
		scalarHandler(cfg.ChatPort, status)
	case "vectorial":
		vectorialHandler(cfg.ChatPort, status)
	default:
		log.Println("TOS not expected")
	}
}

func sequencerHandler(ctx context.Context, addr string, port uint16, verbose bool, status *peer.Status) {
	seqH := seq.NewToSequencerGRPCHandler(addr, port, verbose, status)
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()

		err := seqH.ReceiveMessages(ctx)
		if err != nil {
			log.Printf("Unable to inizialize gRPC server (%v)", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := seqH.SendMessagesToSequencer(ctx)
		if err != nil {
			log.Printf("Unable to send messages anymore (%v)", err)
		}
	}()

	wg.Wait()
}

func scalarHandler(port uint16, status *peer.Status) {
	p2pScalarH := scalar.NewP2PScalarGRPCHandler(port, status)

	errs, _ := errgroup.WithContext(context.Background())

	errs.Go(func() error {
		return p2pScalarH.ReceiveMessages()
	})

	errs.Go(func() error {
		return p2pScalarH.ConnectToPeers()
	})

	errs.Go(func() error {
		p2pScalarH.MultiplexMessages()
		return nil
	})

	if err := errs.Wait(); err != nil {
		log.Fatalf("Something went wrong in grpc connections management (%v)", err)
	}
}

func vectorialHandler(port uint16, status *peer.Status) {
	p2pVectorialH := vectorial.NewP2PVectorialGRPCHandler(port, status)

	errs, _ := errgroup.WithContext(context.Background())

	errs.Go(func() error {
		return p2pVectorialH.ReceiveMessages()
	})

	errs.Go(func() error {
		return p2pVectorialH.ConnectToPeers()
	})

	errs.Go(func() error {
		p2pVectorialH.MultiplexMessages()
		return nil
	})

	errs.Go(func() error {
		return p2pVectorialH.MessageQueueHandler()
	})

	if err := errs.Wait(); err != nil {
		log.Fatalf("Something went wrong in grpc connections management (%v)", err)
	}
}
