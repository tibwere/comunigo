package main

import (
	"fmt"
	"log"

	"gitlab.com/tibwere/comunigo/frontend"
)

func main() {
	ports, err := frontend.GetAvailablePorts("comunigo_peer")
	if err != nil {
		log.Fatalf("Unable to get available peer list")
	}

	fmt.Println(ports)
}
