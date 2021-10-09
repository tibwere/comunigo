package scalar

import (
	"sync"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type P2PScalarGRPCHandler struct {
	*proto.UnimplementedComunigoServer
	comunicationPort  uint16
	peerStatus        *peer.Status
	lockScalar        sync.Mutex
	scalarClock       uint64
	scalarMessagesChs []chan *proto.ScalarClockMessage
	scalarAcksChs     []chan *proto.ScalarClockAck
	pendingMsg        *PendingMessages
	encoder           *protojson.MarshalOptions
}

func NewP2PScalarGRPCHandler(port uint16, status *peer.Status) *P2PScalarGRPCHandler {
	m := &P2PScalarGRPCHandler{
		comunicationPort:  port,
		peerStatus:        status,
		lockScalar:        sync.Mutex{},
		scalarClock:       0,
		scalarMessagesChs: []chan *proto.ScalarClockMessage{},
		scalarAcksChs:     []chan *proto.ScalarClockAck{},
		pendingMsg:        InitPendingMessagesList(status.Members, status.CurrentUsername),
		encoder: &protojson.MarshalOptions{
			Multiline:       false,
			EmitUnpopulated: true,
		},
	}

	for i := 0; i < len(m.peerStatus.Members); i++ {
		if m.peerStatus.Members[i].GetUsername() != m.peerStatus.CurrentUsername {
			m.scalarMessagesChs = append(m.scalarMessagesChs, make(chan *proto.ScalarClockMessage))
			m.scalarAcksChs = append(m.scalarAcksChs, make(chan *proto.ScalarClockAck))
		}
	}
	return m
}
