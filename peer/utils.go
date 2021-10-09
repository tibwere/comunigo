package peer

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

func WaitBeforeSend() {
	delay := rand.Intn(3000)
	log.Printf("Waiting %v millisec ...", delay)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func GetPublicIPAddr() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if !strings.Contains(addr.String(), "127.0.0.1") {
			return strings.Split(addr.String(), "/")[0], nil
		}
	}

	return "", errors.New("no public IP addresses found")
}
