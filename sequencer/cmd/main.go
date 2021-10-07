package main

import (
	"log"
	"sync"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/sequencer/grpchandler"
	"google.golang.org/grpc"
)

func main() {
	var wg sync.WaitGroup
	membersCh := make(chan string)

	config, err := config.SetupSequencer()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	if config.EnableVerbose {
		log.Printf("Start server on port %v\n", config.ChatPort)
	}

	startupServer := &grpchandler.StartupSequencerServer{
		MembersCh: membersCh,
	}

	seqServer := grpchandler.NewSequencerServer(config.ChatPort, config.ChatGroupSize)
	grpcServerToGetPeers := grpc.NewServer()

	wg.Add(2)
	go grpchandler.GetClientsFromRegister(config.RegPort, startupServer, grpcServerToGetPeers, &wg)
	go seqServer.LoadMembers(membersCh, grpcServerToGetPeers, &wg)
	wg.Wait()

	wg.Add(2)
	go seqServer.OrderMessages()
	go grpchandler.ServePeers(seqServer)
	wg.Wait()
}
