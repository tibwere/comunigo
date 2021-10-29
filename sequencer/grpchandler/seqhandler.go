// Package per la gestione della logica applicativa
// basata sullo scambio dei messaggi gRPC con i peer
// nel caso in cui è stata scelta la modalità 'sequencer'
package grpchandler

import (
	"gitlab.com/tibwere/comunigo/proto"
)

// In ottica OO, oggetto che rappresenta
// il sequencer
type SequencerServer struct {
	proto.UnimplementedComunigoServer
	sequenceNumber uint64
	seqCh          chan *proto.RawMessage
	memberCh       chan *proto.PeerInfo
	connections    map[string]chan *proto.SequencerMessage
	port           uint16
	chatGroupSize  uint16
}

// In ottica OO, oggetto che rappresenta
// lo stato ancestrale del sequencer
// che allo startup deve ancora inizializzare
// i propri metadati
type FromRegisterServer struct {
	proto.UnimplementedRegistrationServer
	memberCh chan *proto.PeerInfo
}

// "Costruttore" dell'oggetto FromRegisterServer
func NewFromRegisterServer(memberCh chan *proto.PeerInfo) *FromRegisterServer {
	return &FromRegisterServer{
		memberCh: memberCh,
	}
}

// "Costruttore" dell'oggetto SequencerServer
func NewSequencerServer(port uint16, size uint16, memberCh chan *proto.PeerInfo) *SequencerServer {

	seq := &SequencerServer{
		sequenceNumber: 0,
		seqCh:          make(chan *proto.RawMessage),
		memberCh:       memberCh,
		connections:    make(map[string]chan *proto.SequencerMessage),
		port:           port,
		chatGroupSize:  size,
	}

	return seq
}
