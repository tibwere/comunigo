package p2p

import (
	"log"
	"sort"
	"sync"

	"github.com/go-redis/redis/v8"
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

func InitVectorialMetadata(status *peer.Status, otherNum int) *VectorialMetadata {
	size := BUFFSIZE_FOR_PEER * (len(status.OtherMembers) + 1)
	h := &VectorialMetadata{
		vectorialMessagesChs: []chan *proto.VectorialClockMessage{},
		clock:                []uint64{},
		clockMu:              sync.Mutex{},
		pendingMsg:           []*proto.VectorialClockMessage{},
		receivedCh:           make(chan *proto.VectorialClockMessage, size),
	}

	for i := 0; i < otherNum; i++ {
		h.vectorialMessagesChs = append(h.vectorialMessagesChs, make(chan *proto.VectorialClockMessage, size))
		h.clock = append(h.clock, 0)
	}

	// Serve anche l'entry del processo corrente
	h.clock = append(h.clock, 0)

	// Inizializzazione del clock
	h.initializeClockEntries(status)

	return h
}

func (m *VectorialMetadata) initializeClockEntries(s *peer.Status) {
	var memberUsernames []string
	m.memberIndexs = make(map[string]int, len(s.OtherMembers)+1)

	for _, m := range s.OtherMembers {
		memberUsernames = append(memberUsernames, m.Username)
	}
	memberUsernames = append(memberUsernames, s.CurrentUsername)

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

func (m *VectorialMetadata) SendToAll(mess *proto.VectorialClockMessage, ds *redis.Client, currUser string) error {
	for _, ch := range m.vectorialMessagesChs {
		ch <- mess
	}

	// Questo messaggio può essere direttamente consegnato perché di sicuro
	// rispetta la causalità
	return peer.RPUSHMessage(ds, currUser, mess)
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

func (m *VectorialMetadata) SyncDatastore(ds *redis.Client, currUser string) error {
	deliverables := m.tryToDeliverToDatastore()
	log.Printf("%v new message can be delivered\n", len(deliverables))

	for _, mess := range deliverables {
		log.Printf("Delivering new message: Timestamp: %v, From: %v, Body: '%v'\n", mess.Timestamp, mess.From, mess.Body)

		m.clockMu.Lock()
		m.incrementClockUnlocked(mess.From)
		m.clockMu.Unlock()

		if err := peer.RPUSHMessage(ds, currUser, mess); err != nil {
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

	fromIndex := m.memberIndexs[mess.From]
	for i, messEntry := range mess.Timestamp {
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
