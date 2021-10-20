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
		scalarHandler(ctx, cfg.ChatPort, status)
	case "vectorial":
		vectorialHandler(ctx, cfg.ChatPort, status)
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
			log.Printf("Unable to receive messages anymore (%v)", err)
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

func scalarHandler(ctx context.Context, port uint16, status *peer.Status) {
	p2pScalarH := scalar.NewP2PScalarGRPCHandler(port, status)
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()

		err := p2pScalarH.ReceiveMessages(ctx)
		if err != nil {
			log.Printf("Unable to receive messages anymore (%v)", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := p2pScalarH.ConnectToPeers(ctx)
		if err != nil {
			log.Printf("Sender routines has stopped (%v)", err)
		}
	}()

	go func() {
		defer wg.Done()
		p2pScalarH.MultiplexMessages(ctx)
	}()

	wg.Wait()
}

func vectorialHandler(ctx context.Context, port uint16, status *peer.Status) {
	p2pVectorialH := vectorial.NewP2PVectorialGRPCHandler(port, status)
	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()

		err := p2pVectorialH.ReceiveMessages(ctx)
		if err != nil {
			log.Printf("Unable to receive messages anymore (%v)", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := p2pVectorialH.ConnectToPeers(ctx)
		if err != nil {
			log.Printf("Sender routines has stopped (%v)", err)
		}
	}()

	go func() {
		defer wg.Done()
		p2pVectorialH.MultiplexMessages(ctx)
	}()

	go func() {
		defer wg.Done()

		err := p2pVectorialH.MessageQueueHandler(ctx)
		if err != nil {
			log.Printf("Message queue handler failed (%v)", err)
		}
	}()

	wg.Wait()
}
