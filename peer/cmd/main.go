// Package entry point della logica relativa al peer
// all'interno dell'applicazione comuniGO
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/peer/grpchandler/p2p"
	"gitlab.com/tibwere/comunigo/peer/grpchandler/reg"
	"gitlab.com/tibwere/comunigo/peer/grpchandler/seq"
	"gitlab.com/tibwere/comunigo/peer/webserver"
	"gitlab.com/tibwere/comunigo/utilities"
)

// Launcher del peer
func main() {
	var wg sync.WaitGroup
	ctx := utilities.GetContextForSigHandling()

	rand.Seed(time.Now().UnixNano())

	cfg, err := utilities.InitPeerConfig()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	status, err := peer.Init(cfg.GetRedisAddress())
	if err != nil {
		log.Fatalf("Unable to initialize status (%v)\n", err)
	}

	err = utilities.InitLogger(fmt.Sprintf("peer_%v_main", status.GetExposedIP()))
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	ws := webserver.New(
		cfg.GetWebServerPort(),
		cfg.GetMulticastGroupSize(),
		cfg.GetTOS(),
		cfg.NeedVerbose(),
		status,
	)

	wg.Add(2)
	go ws.Startup(ctx, &wg)
	go internalLogic(ctx, cfg, status, &wg)
	wg.Wait()

	log.Println("Peer is shutting down")
}

// Handler principale per la logica gRPC. Internamente chiama gli altri handkler
// a seconda del tipo di servizio scelto
func internalLogic(ctx context.Context, cfg *utilities.PeerConfig, status *peer.Status, wg *sync.WaitGroup) {
	defer wg.Done()

	regH := reg.NewToRegisterGRPCHandler(cfg.GetRegistrationAddress(), cfg.GetRegistrationPort(), status)
	err := regH.SignToRegister(ctx)
	if err != nil {
		log.Printf("Unable to sign to register node (%v)\n", err)
		return
	}

	switch cfg.GetTOS() {
	case utilities.TOS_CS_SEQUENCER:
		sequencerHandler(ctx, cfg.GetSequencerAddress(), cfg.GetChatPort(), status)
	case utilities.TOS_P2P_SCALAR:
		p2pHandler(ctx, cfg.GetChatPort(), status, p2p.P2P_SCALAR)
	case utilities.TOS_P2P_VECTORIAL:
		p2pHandler(ctx, cfg.GetChatPort(), status, p2p.P2P_VECTORIAL)
	default:
		log.Println("TOS not expected")
	}
}

// Handler gRPC nel caso in cui il tipo di servizio scelto è 'sequencer'
func sequencerHandler(ctx context.Context, addr string, port uint16, status *peer.Status) {
	seqH := seq.NewToSequencerGRPCHandler(addr, port, status)
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

// Handler gRPC nel caso in cui il tipo di servizio scelto è 'scalar' oppure 'vectorial'
func p2pHandler(ctx context.Context, port uint16, status *peer.Status, modality p2p.P2PModality) {
	hnd := p2p.NewP2PHandler(port, status, modality)
	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()

		if err := hnd.ReceiveMessages(ctx); err != nil {
			log.Printf("Unable to receive messages anymore (%v)", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := hnd.ConnectToPeers(ctx); err != nil {
			log.Printf("Sender routines has stopped (%v)", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := hnd.MultiplexMessages(ctx); err != nil {
			log.Printf("Messages multiplexer shutdown (%v)", err)
		}
	}()

	go func() {
		defer wg.Done()

		if err := hnd.MessageQueueHandler(ctx); err != nil {
			log.Printf("Message queue handler failed (%v)", err)
		}
	}()

	wg.Wait()
}
