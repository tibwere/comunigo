package vectorial

import (
	"log"
	"sort"
	"sync"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

type P2PVectorialGRPCHandler struct {
	proto.UnimplementedComunigoServer
	comunicationPort     uint16
	peerStatus           *peer.Status
	vectorialMessagesChs []chan *proto.VectorialClockMessage
	vectorialClock       []uint64
	clockMu              sync.Mutex
	pendingMsg           []*proto.VectorialClockMessage
	memberIndexs         map[string]int
	receivedCh           chan *proto.VectorialClockMessage
}

func NewP2PVectorialGRPCHandler(port uint16, status *peer.Status) *P2PVectorialGRPCHandler {
	h := &P2PVectorialGRPCHandler{
		UnimplementedComunigoServer: proto.UnimplementedComunigoServer{},
		comunicationPort:            port,
		peerStatus:                  status,
		vectorialMessagesChs:        []chan *proto.VectorialClockMessage{},
		vectorialClock:              []uint64{},
		clockMu:                     sync.Mutex{},
		pendingMsg:                  []*proto.VectorialClockMessage{},
		memberIndexs:                initializeClockEntries(status),
		receivedCh:                  make(chan *proto.VectorialClockMessage, 4*(len(status.OtherMembers)+1)),
	}

	for i := 0; i < len(h.peerStatus.OtherMembers); i++ {
		h.vectorialMessagesChs = append(h.vectorialMessagesChs, make(chan *proto.VectorialClockMessage))
	}

	for i := 0; i < len(h.peerStatus.OtherMembers)+1; i++ {
		h.vectorialClock = append(h.vectorialClock, 0)
	}

	return h
}

func initializeClockEntries(s *peer.Status) map[string]int {
	var memberUsernames []string
	indexes := make(map[string]int, len(s.OtherMembers)+1)

	for _, m := range s.OtherMembers {
		memberUsernames = append(memberUsernames, m.Username)
	}
	memberUsernames = append(memberUsernames, s.CurrentUsername)

	sort.Strings(memberUsernames)
	for i, name := range memberUsernames {
		log.Printf("V[%v] -> %v\n", i, name)
		indexes[name] = i
	}

	return indexes
}

func (h *P2PVectorialGRPCHandler) incrementClockUnlocked(member string) {
	index := h.memberIndexs[member]
	h.vectorialClock[index]++
	log.Printf("Incremented V[%v] (entry related to %v). New vectorial clock: %v\n", index, member, h.vectorialClock)
}
