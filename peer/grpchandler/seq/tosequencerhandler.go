// Package per la gestione della logica applicativa
// basata sullo scambio dei messaggi gRPC tra peer e sequencer
// nel caso in cui è stata scelta la modalità 'sequencer'
package seq

import (
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

// In ottica OO, oggetto che rappresenta l'handler della comunicazione
// del peer verso il sequencer
type ToSequencerGRPCHandler struct {
	proto.UnimplementedComunigoServer
	sequencerAddr    string
	comunicationPort uint16
	peerStatus       *peer.Status
}

// "Costruttore" dell'oggetto ToSequencerGRPCHandler
func NewToSequencerGRPCHandler(addr string, port uint16, status *peer.Status) *ToSequencerGRPCHandler {
	return &ToSequencerGRPCHandler{
		sequencerAddr:    addr,
		comunicationPort: port,
		peerStatus:       status,
	}
}
