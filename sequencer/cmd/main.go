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
		log.Printf("Start server on port %v\n", config.SeqPort)
	}

	startupServer := &grpchandler.StartupSequencerServer{
		MembersCh: membersCh,
	}

	seqServer := grpchandler.NewSequencerserver(config.SeqPort, config.ChatGroupSize)
	grpcServer := grpc.NewServer()

	wg.Add(2)
	go grpchandler.GetClientsFromRegister(config.RegPort, startupServer, grpcServer, &wg)
	go seqServer.LoadMembers(membersCh, grpcServer, &wg)
	wg.Wait()

	wg.Add(2)
	go seqServer.OrderMessages(config.SeqPort)
	go grpchandler.ServePeers(seqServer)
	wg.Wait()
}
