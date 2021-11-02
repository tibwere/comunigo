package p2p

import (
	"log"
	"sync"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

// In ottica OO, oggetto che racchiude i metadati basici
// per la comunicazione p2p basata su clock logici scalari
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

// "Costruttore" dell'oggetto ScalarMetadata
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

// "Metodo della classe ScalarMetadata" per inserire un nuovo ack all'interno della
// coda che viene processata dalla goroutine ad-hoc
func (m *ScalarMetadata) InsertNewAck(ack *proto.ScalarClockAck) {
	m.newAckCh <- ack
}

// "Metodo della classe ScalarMetadata" per inserire un nuovo messaggio all'interno della
// coda che viene processata dalla goroutine ad-hoc
func (m *ScalarMetadata) InsertNewMessage(mess *proto.ScalarClockMessage) {
	m.newMessageCh <- mess
}

// "Metodo della classe ScalarMetadata" che permette all'handler della coda pendente
// di accedere al canale dei messaggi in arrivo
func (m *ScalarMetadata) GetMessageCh() <-chan *proto.ScalarClockMessage {
	return m.newMessageCh
}

// "Metodo della classe ScalarMetadata" che permette all'handler della coda pendente
// di accedere al canale degli ack in arrivo
func (m *ScalarMetadata) GetAckCh() <-chan *proto.ScalarClockAck {
	return m.newAckCh
}

// "Metodo della classe ScalarMetadata" che permette di inoltrare a tutti i peer
// un ack e contestualmente autoincrementare il proprio contatore passando
// per la goroutine dedicata al processamento degli ack in ingresso
func (m *ScalarMetadata) AckToAll(ack *proto.ScalarClockAck) {
	for _, ch := range m.scalarAcksChs {
		ch <- ack
	}

	log.Printf("Autoincrement ACK counter after reception of message from %v\n", ack.GetFrom())
	m.newAckCh <- ack

}

// "Metodo della classe ScalarMetadata" che permette di inoltrare a tutti i peer
// un messaggio e contestualmente simulare la ricezione autoincrementando il clock
// e inserendo il messaggio all'interno della lista passando per la goroutine dedicata
func (m *ScalarMetadata) SendToAll(mess *proto.ScalarClockMessage) {
	for _, ch := range m.scalarMessagesChs {
		ch <- mess
	}

	m.clockMu.Lock()
	m.clock++
	m.clockMu.Unlock()

	m.newMessageCh <- mess
}

// "Metodo della classe ScalarMetadata" che permette sia di aggiornare il valore
// del clock corrente a seguito della ricezione di un messaggio rispettando le regole
// previste dall'algoritmo che di generare l'ack da inoltrare ai vari peer connessi
func (m *ScalarMetadata) UpdateClockAtRecv(in *proto.ScalarClockMessage) *proto.ScalarClockAck {
	m.clockMu.Lock()
	// L = max(t, L)
	if m.clock < in.GetTimestamp() {
		m.clock = in.GetTimestamp()
	}

	// L += 1
	m.clock++

	// Generazione del riscontro per il pacchetto ricevuto
	ack := &proto.ScalarClockAck{
		Timestamp: in.GetTimestamp(),
		From:      in.GetFrom(),
	}

	log.Printf("New clock value after update: %v\n", m.clock)
	m.clockMu.Unlock()

	return ack
}

// "Metodo della classe ScalarMetadata" che permette di costruire un nuovo messaggio a partire
// dal "corpo" ricevuto dal frontend
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

// "Metodo della classe ScalarMetadata" che permette alla goroutine dedicata di prelevare
// messaggi da un canale da inoltrare al peer a cui è connesso
func (m *ScalarMetadata) GetIncomingMsgToBeSentCh(index int) <-chan *proto.ScalarClockMessage {
	return m.scalarMessagesChs[index]
}

// "Metodo della classe ScalarMetadata" che permette alla goroutine dedicata di prelevare
// ack da un canale da inoltrare al peer a cui è connesso
func (m *ScalarMetadata) GetIncomingAckToBeSentCh(index int) <-chan *proto.ScalarClockAck {
	return m.scalarAcksChs[index]
}

// "Metodo della classe ScalarMetadata" che permette di inserire un nuovo messaggio ricevuto
// all'interno della coda deu messaggi pendenti ordinandola secondo il valore del clock logico scalare
// e ove necessario dell'identificativo del mittente
func (m *ScalarMetadata) PushIntoPendingList(mess *proto.ScalarClockMessage) {

	inserted := false
	newClock := mess.GetTimestamp()
	newFrom := mess.GetFrom()

	for i, curr := range m.pendingMsg {
		currentClock := curr.GetTimestamp()
		currentFrom := curr.GetFrom()

		if newClock < currentClock || (newClock == currentClock && newFrom < currentFrom) {
			m.pendingMsg = append(m.pendingMsg, &proto.ScalarClockMessage{}) // crea spazio all'interno dello slice per aggiungere un nuovo elemento
			copy(m.pendingMsg[i+1:], m.pendingMsg[i:])                       // shifta verso destra gli elementi dall'i-esimo in poi
			m.pendingMsg[i] = mess                                           // inserisce il messaggio
			inserted = true                                                  // notifica l'inserimento
		}
	}

	if !inserted {
		m.pendingMsg = append(m.pendingMsg, mess)
	}

	// Viene incrementato il presence counter del mittente
	// all'interno della coda per facilitare le operazioni di consegna
	// al livello applicativo
	m.presenceCounter[mess.GetFrom()]++
}

// "Metodo della classe ScalarMetadata" che permette di sincronizzare la lista in memory con
// quella memorizzata nel datastore effettuando il delivery di tutti i messaggi effettivamente
// consegnabili
func (m *ScalarMetadata) SyncDatastore(status *peer.Status) error {
	for _, mess := range m.deliverMessagesIfPossible(status.GetOtherMembers()) {
		log.Printf(
			"Delivered new message (Clock: %v - From: %v)\n",
			mess.GetTimestamp(),
			mess.GetFrom(),
		)
		if err := status.RPUSHMessage(mess); err != nil {
			return err
		}
	}

	return nil
}

// "Metodo della classe ScalarMetadata" che permette di estrarre la sottolista dei messaggi consegnabili
// dalla lista dei messaggi pendenti
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

// "Metodo della classe ScalarMetadata" che verifica se per ogni altro peer connesso
// sono presenti effettivamente dei messaggi in coda
//
// (una delle condizioni di consegnabilità dell'algoritmo lo prevede)
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

// "Metodo della classe ScalarMetadata" che permette di incrementare il contatore degli ack ricevuti
//
// (necessario per gestire la condizione sugli ack della consegnabilità)
func (m *ScalarMetadata) IncrementAckCounter(ack *proto.ScalarClockAck) {
	m.receivedAcks[ack.String()]++
}
