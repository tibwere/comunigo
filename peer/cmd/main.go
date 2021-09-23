package main

import (
	"fmt"
	"sync"

	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/peer"
)

func main() {
	var wg sync.WaitGroup
	usernameCh := make(chan string)

	config, err := config.SetupPeerConfiguration()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Dimensione: %v", config.ChatGroupSize)

	wg.Add(2)
	go peer.StartupWebServer(config.WebServerPort, config.ChatGroupSize, usernameCh, &wg)
	go peer.SignToRegister(config.RegHostname, config.RegPort, usernameCh)
	wg.Wait()
}
