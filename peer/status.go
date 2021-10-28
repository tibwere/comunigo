package peer

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/go-redis/redis/v8"
	"gitlab.com/tibwere/comunigo/proto"
)

type Status struct {
	currentUsername string
	otherMembers    []*proto.PeerInfo
	datastore       *redis.Client
	frontBackCh     chan string
	exposedIP       string
}

func Init(redisAddr string) (*Status, error) {
	ip, err := retrieveIP()
	if err != nil {
		return nil, err
	} else {
		s := &Status{
			currentUsername: "",
			otherMembers:    []*proto.PeerInfo{},
			frontBackCh:     make(chan string),
			exposedIP:       ip,
		}

		s.initDatastore(redisAddr)
		return s, nil
	}
}

func (s *Status) GetExposedIP() string {
	return s.exposedIP
}

func (s *Status) GetOtherMembers() []*proto.PeerInfo {
	return s.otherMembers
}

func (s *Status) GetSpecificMember(index int) *proto.PeerInfo {
	return s.otherMembers[index]
}

func (s *Status) SetUsername(username string) error {
	if s.currentUsername != "" {
		return fmt.Errorf("unable to set twice username")
	} else {
		s.currentUsername = username
		return nil
	}
}

func (s *Status) GetCurrentUsername() string {
	return s.currentUsername
}

func (s *Status) NotYetSigned() bool {
	return s.currentUsername == ""
}

func (s *Status) GetFromFrontendBackendChannel() <-chan string {
	return s.frontBackCh
}

func (s *Status) PushIntoFrontendBackendChannel(message string) {
	s.frontBackCh <- message
}

func (s *Status) InsertNewMember(newMember *proto.PeerInfo) {
	s.otherMembers = append(s.otherMembers, newMember)
}

func retrieveIP() (string, error) {
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
