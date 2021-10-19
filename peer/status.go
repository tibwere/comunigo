package peer

import (
	"github.com/go-redis/redis/v8"
	"gitlab.com/tibwere/comunigo/proto"
)

type Status struct {
	CurrentUsername string
	OtherMembers    []*proto.PeerInfo
	Datastore       *redis.Client
	FrontBackCh     chan string
	PublicIP        string
}

func Init(redisAddr string) (*Status, error) {
	ip, err := GetPublicIPAddr()
	if err != nil {
		return nil, err
	} else {
		return &Status{
			CurrentUsername: "",
			OtherMembers:    []*proto.PeerInfo{},
			Datastore:       InitDatastore(redisAddr),
			FrontBackCh:     make(chan string),
			PublicIP:        ip,
		}, nil
	}
}
