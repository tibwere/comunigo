package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/peer/grpchandler"
	"gitlab.com/tibwere/comunigo/peer/webserver"
	"golang.org/x/sync/errgroup"
)

func main() {
	var wg sync.WaitGroup

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

	webserver := webserver.New(cfg.WebServerPort, cfg.ChatGroupSize, status)
	grpcHandler := grpchandler.New(cfg, status)

	wg.Add(2)
	go webserver.Startup(&wg)
	go func() {
		defer wg.Done()

		err = grpcHandler.SignToRegister()
		if err != nil {
			log.Fatalf("Unable to sign to register node")
		}

		errs, _ := errgroup.WithContext(context.Background())
		errs.Go(func() error {
			return grpcHandler.ReceiveMessages()
		})
		errs.Go(func() error {
			return grpcHandler.SendMessages()
		})

		if err = errs.Wait(); err != nil {
			log.Fatalf("Something went wrong in grpc connections management (%v)", err)
		}
	}()
	wg.Wait()
}
