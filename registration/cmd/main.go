package main

import (
	"log"
	"sync"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/registration"
)

func main() {
	var wg sync.WaitGroup

	config, err := config.SetupRegistrationServerConfiguration()
	if err != nil {
		log.Fatalf("Unable to load configurations (%v)\n", err)
	}

	if config.EnableVerbose {
		log.Printf("In ascolto sulla porta 2929 (Dimensione del gruppo: %v)\n", config.ChatGroupSize)
	}

	regServer := registration.NewRegistrationServer(config.ChatGroupSize)

	wg.Add(2)
	go regServer.UpdateMembers()
	go registration.ServeSignRequests(config.RegPort, regServer)
	wg.Wait()
}
