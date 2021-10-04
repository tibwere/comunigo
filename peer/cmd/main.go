package main

import (
	"log"
	"sync"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/peer/grpchandler"
	"gitlab.com/tibwere/comunigo/peer/webserver"
)

func main() {
	var wg sync.WaitGroup

	config, err := config.SetupPeer()
	if err != nil {
		panic(err)
	}

	status := peer.Init(config.RedisHostname)
	webserver := webserver.New(config.WebServerPort, config.ChatGroupSize, status)
	grpcHandler := grpchandler.New(config, status)

	wg.Add(2)
	go webserver.Startup(&wg)
	go func() {
		defer wg.Done()

		err = grpcHandler.SignToRegister()
		if err != nil {
			log.Fatalf("Unable to sign to register node")
		}

		var childWg sync.WaitGroup
		childWg.Add(2)
		go grpcHandler.ReceiveMessages(&childWg)
		go grpcHandler.SendMessages(&childWg)
		childWg.Wait()
	}()
	wg.Wait()
}
