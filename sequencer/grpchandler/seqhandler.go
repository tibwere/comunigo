package grpchandler

import (
	"gitlab.com/tibwere/comunigo/proto"
)

type SequencerServer struct {
	proto.UnimplementedComunigoServer
	sequenceNumber uint64
	seqCh          chan *proto.RawMessage
	memberCh       chan *proto.PeerInfo
	connections    map[string]chan *proto.SequencerMessage
	port           uint16
	chatGroupSize  uint16
}

type FromRegisterServer struct {
	proto.UnimplementedRegistrationServer
	memberCh chan *proto.PeerInfo
}

func NewFromRegisterServer(memberCh chan *proto.PeerInfo) *FromRegisterServer {
	return &FromRegisterServer{
		memberCh: memberCh,
	}
}

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
