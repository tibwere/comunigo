package peer

import (
	"errors"
	"net"
	"strings"

	"github.com/go-redis/redis/v8"
	"gitlab.com/tibwere/comunigo/proto"
)

type Status struct {
	CurrentUsername string
	Members         []*proto.ClientInfo
	Datastore       *redis.Client
	UsernameCh      chan string
	InvalidCh       chan bool
	DoneCh          chan bool
	RawMessageCh    chan string
	PublicIP        string
}

func Init(redisAddr string) (*Status, error) {
	ip, err := getPublicIPAddr()
	if err != nil {
		return nil, err
	} else {
		return &Status{
			CurrentUsername: "",
			Members:         []*proto.ClientInfo{},
			Datastore:       InitDatastore(redisAddr),
			UsernameCh:      make(chan string),
			InvalidCh:       make(chan bool),
			DoneCh:          make(chan bool),
			RawMessageCh:    make(chan string),
			PublicIP:        ip,
		}, nil
	}
}

func getPublicIPAddr() (string, error) {
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
