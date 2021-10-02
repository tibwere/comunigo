package main

import (
	"log"
	"sync"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/peer/grpchandler"
)

var currentUser string

func main() {
	var wg sync.WaitGroup

	config, err := config.SetupPeerConfig()
	if err != nil {
		panic(err)
	}

	channels := peer.InitChannels()
	currentUser = ""

	wg.Add(2)
	go peer.StartupWebServer(
		config.WebServerPort,
		config.ChatGroupSize,
		channels,
		&currentUser,
		&wg,
	)
	go func() {
		defer wg.Done()

		currentUser, err = grpchandler.SignToRegister(config.RegHostname, config.RegPort, channels.UsernameCh)
		if err != nil {
			log.Fatalf("Unable to sign to register node")
		}

		var childWg sync.WaitGroup
		childWg.Add(2)
		go grpchandler.ReceiveMessages(config.SeqPort, &childWg)
		go grpchandler.SendMessages(config.SeqHostname, config.SeqPort, currentUser, channels.RawMessageCh, &childWg)
		childWg.Wait()
	}()
	wg.Wait()
}
