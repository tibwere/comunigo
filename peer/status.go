package peer

import (
	"github.com/go-redis/redis/v8"
	"gitlab.com/tibwere/comunigo/proto"
)

type Status struct {
	CurrentUsername string
	Members         []*proto.PeerInfo
	Datastore       *redis.Client
	UsernameCh      chan string
	InvalidCh       chan bool
	DoneCh          chan bool
	RawMessageCh    chan string
	PublicIP        string
}

func Init(redisAddr string) (*Status, error) {
	ip, err := GetPublicIPAddr()
	if err != nil {
		return nil, err
	} else {
		return &Status{
			CurrentUsername: "",
			Members:         []*proto.PeerInfo{},
			Datastore:       InitDatastore(redisAddr),
			UsernameCh:      make(chan string),
			InvalidCh:       make(chan bool),
			DoneCh:          make(chan bool),
			RawMessageCh:    make(chan string),
			PublicIP:        ip,
		}, nil
	}
}
