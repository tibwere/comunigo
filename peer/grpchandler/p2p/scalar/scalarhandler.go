package scalar

import (
	"sync"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

type P2PScalarGRPCHandler struct {
	proto.UnimplementedComunigoServer
	comunicationPort  uint16
	peerStatus        *peer.Status
	lockScalar        sync.Mutex
	scalarClock       uint64
	scalarMessagesChs []chan *proto.ScalarClockMessage
	scalarAcksChs     []chan *proto.ScalarClockAck
	pendingMsg        *PendingMessages
}

func NewP2PScalarGRPCHandler(port uint16, status *peer.Status) *P2PScalarGRPCHandler {
	h := &P2PScalarGRPCHandler{
		comunicationPort:  port,
		peerStatus:        status,
		lockScalar:        sync.Mutex{},
		scalarClock:       0,
		scalarMessagesChs: []chan *proto.ScalarClockMessage{},
		scalarAcksChs:     []chan *proto.ScalarClockAck{},
		pendingMsg:        InitPendingMessagesList(status.OtherMembers, status.CurrentUsername),
	}

	for i := 0; i < len(h.peerStatus.OtherMembers); i++ {
		h.scalarMessagesChs = append(h.scalarMessagesChs, make(chan *proto.ScalarClockMessage))
		h.scalarAcksChs = append(h.scalarAcksChs, make(chan *proto.ScalarClockAck))
	}
	return h
}
