package p2p

import (
	"log"
	"sort"
	"sync"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

type ScalarMetadata struct {
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

func InitScalarMetadata(members []*proto.PeerInfo) *ScalarMetadata {
	size := len(members) * BUFFSIZE_FOR_PEER
	h := &ScalarMetadata{
		clockMu:           sync.Mutex{},
		clock:             0,
		scalarMessagesChs: []chan *proto.ScalarClockMessage{},
		scalarAcksChs:     []chan *proto.ScalarClockAck{},
		newMessageCh:      make(chan *proto.ScalarClockMessage, size),
		newAckCh:          make(chan *proto.ScalarClockAck, size),
		pendingMsg:        []*proto.ScalarClockMessage{},
		presenceCounter:   make(map[string]int),
		receivedAcks:      make(map[string]int),
	}

	for _, m := range members {
		h.scalarMessagesChs = append(h.scalarMessagesChs, make(chan *proto.ScalarClockMessage, size))
		h.scalarAcksChs = append(h.scalarAcksChs, make(chan *proto.ScalarClockAck, size))
		h.presenceCounter[m.GetUsername()] = 0
	}

	return h
}

func (m *ScalarMetadata) InsertNewAck(ack *proto.ScalarClockAck) {
	m.newAckCh <- ack
}

func (m *ScalarMetadata) InsertNewMessage(mess *proto.ScalarClockMessage) {
	m.newMessageCh <- mess
}

func (m *ScalarMetadata) GetMessageCh() <-chan *proto.ScalarClockMessage {
	return m.newMessageCh
}

func (m *ScalarMetadata) GetAckCh() <-chan *proto.ScalarClockAck {
	return m.newAckCh
}

func (m *ScalarMetadata) AckToAll(ack *proto.ScalarClockAck, autoIncrement bool) {
	for _, ch := range m.scalarAcksChs {
		ch <- ack
	}

	if autoIncrement {
		log.Printf("Autoincrement ACK counter after reception of message from %v\n", ack.GetFrom())
		m.newAckCh <- ack
	}
}

func (m *ScalarMetadata) SendToAll(mess *proto.ScalarClockMessage) {
	for _, ch := range m.scalarMessagesChs {
		ch <- mess
	}

	// Simulazione di ricezione
	m.clockMu.Lock()
	m.clock++
	m.clockMu.Unlock()

	m.newMessageCh <- mess
}

func (m *ScalarMetadata) UpdateClockAtRecv(in *proto.ScalarClockMessage) *proto.ScalarClockAck {
	m.clockMu.Lock()
	// L = max(t, L)
	if m.clock < in.GetTimestamp() {
		m.clock = in.GetTimestamp()
	}

	// L += 1
	m.clock++

	// Invio del riscontro per il pacchetto ricevuto
	ack := &proto.ScalarClockAck{
		Timestamp: in.GetTimestamp(),
		From:      in.GetFrom(),
	}

	log.Printf("New clock value after update: %v\n", m.clock)
	m.clockMu.Unlock()

	return ack
}

func (m *ScalarMetadata) GenerateNewMessage(from string, body string) *proto.ScalarClockMessage {
	m.clockMu.Lock()
	m.clock++
	newMessage := &proto.ScalarClockMessage{
		Timestamp: m.clock,
		From:      from,
		Body:      body,
	}
	m.clockMu.Unlock()

	return newMessage
}

func (m *ScalarMetadata) GetIncomingMsgToBeSentCh(index int) <-chan *proto.ScalarClockMessage {
	return m.scalarMessagesChs[index]
}

func (m *ScalarMetadata) GetIncomingAckToBeSentCh(index int) <-chan *proto.ScalarClockAck {
	return m.scalarAcksChs[index]
}

func (m *ScalarMetadata) PushIntoPendingList(mess *proto.ScalarClockMessage) {
	m.pendingMsg = append(m.pendingMsg, mess)

	sort.Slice(m.pendingMsg, func(i, j int) bool {

		iClock := m.pendingMsg[i].GetTimestamp()
		jClock := m.pendingMsg[j].GetTimestamp()
		iFrom := m.pendingMsg[i].GetFrom()
		jFrom := m.pendingMsg[j].GetFrom()

		return iClock < jClock || (iClock == jClock && iFrom < jFrom)
	})
	m.presenceCounter[mess.GetFrom()]++
}

func (m *ScalarMetadata) SyncDatastore(status *peer.Status) error {
	for _, mess := range m.deliverMessagesIfPossible(status.GetOtherMembers()) {
		log.Printf("Delivered new message (Clock: %v - From: %v)\n", mess.GetTimestamp(), mess.GetFrom())
		if err := status.RPUSHMessage(mess); err != nil {
			return err
		}
	}

	return nil
}

func (m *ScalarMetadata) deliverMessagesIfPossible(others []*proto.PeerInfo) []*proto.ScalarClockMessage {
	var deliverList []*proto.ScalarClockMessage

	if len(m.pendingMsg) == 0 {
		return deliverList
	}

	firstMsg := m.pendingMsg[0]
	firstAck := &proto.ScalarClockAck{
		Timestamp: firstMsg.GetTimestamp(),
		From:      firstMsg.GetFrom(),
	}
	nMember := len(others)

	log.Printf("Received %v/%v acks for [%v]\n", m.receivedAcks[firstAck.String()], nMember, firstMsg)

	if m.receivedAcks[firstAck.String()] == nMember && m.thereAreMessagesFromAllInQueue(firstMsg.GetFrom()) {
		deliverList = append(deliverList, firstMsg)
		m.pendingMsg = m.pendingMsg[1:]
		m.presenceCounter[firstAck.GetFrom()]--
	}

	return deliverList
}

func (m *ScalarMetadata) thereAreMessagesFromAllInQueue(actualFrom string) bool {
	for from, presences := range m.presenceCounter {
		if from != actualFrom && presences == 0 {
			log.Printf("Member %v does not have messages in queue at this moment. Cannot deliver message\n", from)
			return false
		}
	}

	log.Println("All other members has at least a message in queue at this moment.")

	return true
}

func (m *ScalarMetadata) IncrementAckCounter(ack *proto.ScalarClockAck) {
	m.receivedAcks[ack.String()]++
}
