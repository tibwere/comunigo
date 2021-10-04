package grpchandler

import (
	"gitlab.com/tibwere/comunigo/config"
	"gitlab.com/tibwere/comunigo/peer"
)

type GrpcHandler struct {
	registerAddr  string
	registerPort  uint16
	sequencerAddr string
	sequencerPort uint16
	peerStatus    *peer.Status
}

func New(config *config.PeerConfig, status *peer.Status) *GrpcHandler {
	return &GrpcHandler{
		registerAddr:  config.RegHostname,
		registerPort:  config.RegPort,
		sequencerAddr: config.SeqHostname,
		sequencerPort: config.SeqPort,
		peerStatus:    status,
	}
}
