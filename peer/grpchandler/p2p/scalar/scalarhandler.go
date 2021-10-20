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
	clockMu           sync.Mutex
	clock             uint64
	scalarMessagesChs []chan *proto.ScalarClockMessage
	scalarAcksChs     []chan *proto.ScalarClockAck
	newMessageCh      chan *proto.ScalarClockMessage
	newAckCh          chan *proto.ScalarClockAck
	pendingMsg        []*proto.ScalarClockMessage
	presenceCounter   map[string]int
	receivedAcks      map[string]int
}

func NewP2PScalarGRPCHandler(port uint16, status *peer.Status) *P2PScalarGRPCHandler {
	h := &P2PScalarGRPCHandler{
		UnimplementedComunigoServer: proto.UnimplementedComunigoServer{},
		comunicationPort:            port,
		peerStatus:                  status,
		clockMu:                     sync.Mutex{},
		clock:                       0,
		scalarMessagesChs:           []chan *proto.ScalarClockMessage{},
		scalarAcksChs:               []chan *proto.ScalarClockAck{},
		newMessageCh:                make(chan *proto.ScalarClockMessage),
		newAckCh:                    make(chan *proto.ScalarClockAck),
		pendingMsg:                  []*proto.ScalarClockMessage{},
		presenceCounter:             make(map[string]int),
		receivedAcks:                make(map[string]int),
	}

	for _, m := range h.peerStatus.OtherMembers {
		h.scalarMessagesChs = append(h.scalarMessagesChs, make(chan *proto.ScalarClockMessage))
		h.scalarAcksChs = append(h.scalarAcksChs, make(chan *proto.ScalarClockAck))
		h.presenceCounter[m.Username] = 0
	}

	return h
}
