package p2p

import (
	"log"
	"sort"
	"sync"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

type VectorialMetadata struct {
	vectorialMessagesChs []chan *proto.VectorialClockMessage
	clock                []uint64
	clockMu              sync.Mutex
	pendingMsg           []*proto.VectorialClockMessage
	memberIndexs         map[string]int
	receivedCh           chan *proto.VectorialClockMessage
}

func InitVectorialMetadata(current string, others []*proto.PeerInfo) *VectorialMetadata {
	size := BUFFSIZE_FOR_PEER * (len(others) + 1)
	h := &VectorialMetadata{
		vectorialMessagesChs: []chan *proto.VectorialClockMessage{},
		clock:                []uint64{},
		clockMu:              sync.Mutex{},
		pendingMsg:           []*proto.VectorialClockMessage{},
		receivedCh:           make(chan *proto.VectorialClockMessage, size),
	}

	for i := 0; i < len(others); i++ {
		h.vectorialMessagesChs = append(h.vectorialMessagesChs, make(chan *proto.VectorialClockMessage, size))
		h.clock = append(h.clock, 0)
	}

	// Serve anche l'entry del processo corrente
	h.clock = append(h.clock, 0)

	// Inizializzazione del clock
	h.initializeClockEntries(current, others)

	return h
}

func (m *VectorialMetadata) initializeClockEntries(current string, others []*proto.PeerInfo) {
	var memberUsernames []string
	m.memberIndexs = make(map[string]int, len(others)+1)

	for _, m := range others {
		memberUsernames = append(memberUsernames, m.GetUsername())
	}
	memberUsernames = append(memberUsernames, current)

	sort.Strings(memberUsernames)
	for i, name := range memberUsernames {
		log.Printf("V[%v] -> %v\n", i, name)
		m.memberIndexs[name] = i
	}
}

func (m *VectorialMetadata) GenerateNewMessage(from string, body string) *proto.VectorialClockMessage {
	m.clockMu.Lock()
	m.incrementClockUnlocked(from)
	newMessage := &proto.VectorialClockMessage{
		Timestamp: m.clock,
		From:      from,
		Body:      body,
	}
	m.clockMu.Unlock()

	return newMessage
}

func (m *VectorialMetadata) InsertNewMessage(mess *proto.VectorialClockMessage) {
	m.receivedCh <- mess
}

func (m *VectorialMetadata) incrementClockUnlocked(member string) {
	index := m.memberIndexs[member]
	m.clock[index]++
	log.Printf("Incremented V[%v] (entry related to %v). New vectorial clock: %v\n", index, member, m.clock)
}

func (m *VectorialMetadata) SendToAll(mess *proto.VectorialClockMessage) {
	for _, ch := range m.vectorialMessagesChs {
		ch <- mess
	}
}

func (m *VectorialMetadata) GetIncomingMsgToBeSentCh(index int) <-chan *proto.VectorialClockMessage {
	return m.vectorialMessagesChs[index]
}

func (m *VectorialMetadata) GetReceivedCh() <-chan *proto.VectorialClockMessage {
	return m.receivedCh
}

func (m *VectorialMetadata) PushIntoPendingList(mess *proto.VectorialClockMessage) {
	m.pendingMsg = append(m.pendingMsg, mess)
}

func (m *VectorialMetadata) SyncDatastore(status *peer.Status) error {
	deliverables := m.tryToDeliverToDatastore()
	log.Printf("%v new message can be delivered\n", len(deliverables))

	for _, mess := range deliverables {
		log.Printf(
			"Delivering new message: Timestamp: %v, From: %v, Body: '%v'\n",
			mess.GetTimestamp(),
			mess.GetFrom(),
			mess.GetBody(),
		)

		m.clockMu.Lock()
		m.incrementClockUnlocked(mess.GetFrom())
		m.clockMu.Unlock()

		if err := status.RPUSHMessage(mess); err != nil {
			return err
		}
	}

	return nil
}

func (m *VectorialMetadata) tryToDeliverToDatastore() []*proto.VectorialClockMessage {

	var deliverables []*proto.VectorialClockMessage

	for _, mess := range m.pendingMsg {
		if m.isDeliverable(mess) {
			deliverables = append(deliverables, mess)
		}
	}

	return deliverables
}

func (m *VectorialMetadata) isDeliverable(mess *proto.VectorialClockMessage) bool {

	fromIndex := m.memberIndexs[mess.GetFrom()]
	for i, messEntry := range mess.GetTimestamp() {
		localEntry := m.clock[i]
		if i == fromIndex && messEntry != localEntry+1 {
			return false
		}

		if i != fromIndex && messEntry > localEntry {
			return false
		}
	}

	return true
}
