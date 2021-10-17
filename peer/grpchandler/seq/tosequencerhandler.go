package seq

import (
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

type ToSequencerGRPCHandler struct {
	proto.UnimplementedComunigoServer
	sequencerAddr    string
	comunicationPort uint16
	peerStatus       *peer.Status
	verbose          bool
}

func NewToSequencerGRPCHandler(addr string, port uint16, verbose bool, status *peer.Status) *ToSequencerGRPCHandler {
	return &ToSequencerGRPCHandler{
		sequencerAddr:    addr,
		comunicationPort: port,
		peerStatus:       status,
		verbose:          verbose,
	}
}
