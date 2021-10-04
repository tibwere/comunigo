package peer

import (
	"github.com/go-redis/redis"
	"gitlab.com/tibwere/comunigo/proto"
)

type Status struct {
	CurrentUsername string
	Members         []*proto.ClientInfo
	Datastore       *redis.Client
	UsernameCh      chan string
	MemberCh        chan *proto.ClientInfo
	InvalidCh       chan bool
	DoneCh          chan bool
	RawMessageCh    chan string
}

func Init(redisAddr string) *Status {
	return &Status{
		CurrentUsername: "",
		Members:         []*proto.ClientInfo{},
		Datastore:       InitDatastore(redisAddr),
		UsernameCh:      make(chan string),
		MemberCh:        make(chan *proto.ClientInfo),
		InvalidCh:       make(chan bool),
		DoneCh:          make(chan bool),
		RawMessageCh:    make(chan string),
	}
}
