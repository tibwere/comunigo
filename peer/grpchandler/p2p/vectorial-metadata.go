package p2p

import (
	"log"
	"sort"
	"sync"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

// In ottica OO, oggetto che racchiude i metadati basici
// per la comunicazione p2p basata su clock logici vettoriale
type VectorialMetadata struct {
	vectorialMessagesChs []chan *proto.VectorialClockMessage
	clock                []uint64
	clockMu              sync.Mutex
	pendingMsg           []*proto.VectorialClockMessage
	memberIndexs         map[string]int
	receivedCh           chan *proto.VectorialClockMessage
}

// "Costruttore" dell'oggetto VectorialMetadata
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

	// Inizializzazione dell'entry del processo corrente
	h.clock = append(h.clock, 0)

	h.mapClockEntriesToUsernames(current, others)

	return h
}

// "Metodo della classe VectorialMetadata" che permette di inizializzare la mappa della classe
// adibita al mapping fra gli username dei peer connessi al gruppo di multicast
// e le entry del clock logico vettoriale
func (m *VectorialMetadata) mapClockEntriesToUsernames(current string, others []*proto.PeerInfo) {
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

// "Metodo della classe VectorialMetadata" che permette di costruire un nuovo messaggio a partire
// dal "corpo" ricevuto dal frontend
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

// "Metodo della classe VectorialMetadata" per inserire un nuovo messaggio all'interno della
// coda che viene processata dalla goroutine ad-hoc
func (m *VectorialMetadata) InsertNewMessage(mess *proto.VectorialClockMessage) {
	m.receivedCh <- mess
}

// "Metodo della classe VectorialMetadata" per l'incremento del clock
//
// n.b. Questo metodo dev'essere invocato soltanto dopo aver preso un lock
func (m *VectorialMetadata) incrementClockUnlocked(member string) {
	index := m.memberIndexs[member]
	m.clock[index]++
	log.Printf("Incremented V[%v] (entry related to %v). New vectorial clock: %v\n", index, member, m.clock)
}

// "Metodo della classe VectorialMetadata" che permette di inoltrare
// a tutti i peer un messaggio
func (m *VectorialMetadata) SendToAll(mess *proto.VectorialClockMessage) {
	for _, ch := range m.vectorialMessagesChs {
		ch <- mess
	}
}

// "Metodo della classe VectorialMetadata" che permette alla goroutine dedicata di prelevare
// messaggi da un canale da inoltrare al peer a cui è connesso
func (m *VectorialMetadata) GetIncomingMsgToBeSentCh(index int) <-chan *proto.VectorialClockMessage {
	return m.vectorialMessagesChs[index]
}

// "Metodo della classe VectorialMetadata" che permette alla goroutine dedicata di prelevare
// i messaggi ricevuti da inserire all'interno della lista dei pendenti
func (m *VectorialMetadata) GetReceivedCh() <-chan *proto.VectorialClockMessage {
	return m.receivedCh
}

// "Metodo della classe VectorialMetadata" che permette di inserire un messaggio all'interno della lista
// dei messaggi pendenti
func (m *VectorialMetadata) PushIntoPendingList(mess *proto.VectorialClockMessage) {
	m.pendingMsg = append(m.pendingMsg, mess)
}

// "Metodo della classe VectorialMetadata" che permette di sincronizzare la lista in memory con
// quella memorizzata nel datastore effettuando il delivery di tutti i messaggi effettivamente
// consegnabili
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

// "Metodo della classe VectorialMetadata" che permette di estrarre la sottolista dei messaggi consegnabili
// dalla lista dei messaggi pendenti
func (m *VectorialMetadata) tryToDeliverToDatastore() []*proto.VectorialClockMessage {

	var deliverables []*proto.VectorialClockMessage

	for _, mess := range m.pendingMsg {
		if m.isDeliverable(mess) {
			deliverables = append(deliverables, mess)
		}
	}

	return deliverables
}

// "Metodo della classe VectorialMetadata" che permette di verificare se un messaggio
// è o meno consegnabile
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
