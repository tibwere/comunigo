package grpchandler

import (
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

type SequencerServer struct {
	proto.UnimplementedComunigoServer
	proto.UnimplementedRegistrationServer

	// Questo grpc server è salvato tra i metadati basici
	// del sequencers server perché deve essere stoppato
	// nel momento in cui non sono stati ricevuti tutti i
	// peer dal servizio di registrazione
	fromRegToSeqGRPCserver *grpc.Server
	sequenceNumber         uint64
	seqCh                  chan *proto.RawMessage
	connections            map[string]chan *proto.SequencerMessage
	port                   uint16
	chatGroupSize          uint16
}

func NewSequencerServer(port uint16, size uint16) *SequencerServer {

	seq := &SequencerServer{
		fromRegToSeqGRPCserver: grpc.NewServer(),
		sequenceNumber:         0,
		seqCh:                  make(chan *proto.RawMessage),
		connections:            make(map[string]chan *proto.SequencerMessage),
		port:                   port,
		chatGroupSize:          size,
	}

	return seq
}
