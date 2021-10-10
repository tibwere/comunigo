package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/peer/grpchandler/p2p/scalar"
	"gitlab.com/tibwere/comunigo/peer/grpchandler/reg"
	"gitlab.com/tibwere/comunigo/peer/grpchandler/seq"
	"gitlab.com/tibwere/comunigo/peer/webserver"
	"golang.org/x/sync/errgroup"
)

func main() {
	var wg sync.WaitGroup

	rand.Seed(time.Now().UnixNano())

	cfg, err := config.SetupPeer()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	status, err := peer.Init(cfg.RedisHostname)
	if err != nil {
		log.Fatalf("Unable to initialize status (%v)\n", err)
	}

	err = config.InitLogger(fmt.Sprintf("peer_%v_main", status.PublicIP))
	if err != nil {
		log.Fatalf("Unable to setup log file (%v)\n", err)
	}

	wg.Add(2)
	ws := webserver.New(cfg.WebServerPort, cfg.ChatGroupSize, status)
	go ws.Startup(&wg)
	go func() {
		defer wg.Done()

		regH := reg.NewToRegisterGRPCHandler(cfg.RegHostname, cfg.RegPort, status)
		err = regH.SignToRegister()
		if err != nil {
			log.Fatalf("Unable to sign to register node")
		}

		switch cfg.TypeOfService {
		case "sequencer":
			sequencerHandler(cfg.SeqHostname, cfg.ChatPort, status)
		case "scalar":
			scalarHandler(cfg.ChatPort, status)
		case "vectorial":
			vectorialHandler()
		default:
			log.Fatalf("TOS not expected")
		}

	}()
	wg.Wait()
}

func sequencerHandler(addr string, port uint16, status *peer.Status) {
	seqH := seq.NewToSequencerGRPCHandler(addr, port, status)
	errs, _ := errgroup.WithContext(context.Background())

	errs.Go(func() error {
		return seqH.ReceiveMessages()
	})
	errs.Go(func() error {
		return seqH.SendMessagesToSequencer()
	})

	if err := errs.Wait(); err != nil {
		log.Fatalf("Something went wrong in grpc connections management (%v)", err)
	}
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

func vectorialHandler() {
	log.Fatalf("Not implemented yet")
}
