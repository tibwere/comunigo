package reg

import "gitlab.com/tibwere/comunigo/peer"

type ToRegisterGRPCHandler struct {
	registerAddr string
	registerPort uint16
	peerStatus   *peer.Status
}

func NewToRegisterGRPCHandler(addr string, port uint16, status *peer.Status) *ToRegisterGRPCHandler {
	return &ToRegisterGRPCHandler{
		registerAddr: addr,
		registerPort: port,
		peerStatus:   status,
	}
}
