package scalar

import (
	"log"
	"sort"
	"sync"

	"gitlab.com/tibwere/comunigo/proto"
)

type PendingMessages struct {
	lock            *sync.Mutex
	queue           []*proto.ScalarClockMessage
	receivedAcks    map[string]int
	presenceCounter map[string]int
	otherFroms      []string
}

func InitPendingMessagesList(allMembers []*proto.PeerInfo, currUser string) *PendingMessages {
	pm := &PendingMessages{
		lock:            &sync.Mutex{},
		queue:           []*proto.ScalarClockMessage{},
		receivedAcks:    map[string]int{},
		presenceCounter: make(map[string]int),
		otherFroms:      []string{},
	}

	for _, m := range allMembers {
		if m.Username != currUser {
			pm.otherFroms = append(pm.otherFroms, m.Username)
			pm.presenceCounter[m.Username] = 0
		}
	}

	return pm
}

func (pm *PendingMessages) Insert(newMessage *proto.ScalarClockMessage) {
	pm.lock.Lock()
	log.Printf("Insert [%v] into pendant queue\n", newMessage)
	pm.queue = append(pm.queue, newMessage)
	pm.presenceCounter[newMessage.From]++
	sort.Slice(pm.queue, func(i, j int) bool {

		iClock := pm.queue[i].ScalarClock
		jClock := pm.queue[j].ScalarClock
		iFrom := pm.queue[i].From
		jFrom := pm.queue[i].From

		return iClock < jClock || (iClock == jClock && iFrom < jFrom)
	})
	pm.lock.Unlock()
}

func isAckedMessage(mess *proto.ScalarClockMessage, ack *proto.ScalarClockAck) bool {
	return mess.From == ack.From && mess.ScalarClock == ack.ScalarClock
}

func (pm *PendingMessages) thereAreMessagesFromAllInQueue(actualFrom string) bool {
	for from, presences := range pm.presenceCounter {
		if from != actualFrom && presences == 0 {
			log.Printf("Member %v does not have messages in queue at this moment. Cannot deliver message\n", from)
			return false
		}
	}

	log.Println("All other members has at least a message in queue at this moment.")

	return true
}

func (pm *PendingMessages) CheckIfIsReadyToDelivered(currUser string) *proto.ScalarClockMessage {
	var deliverMsg *proto.ScalarClockMessage
	var expectedAcks int
	canDeliver := false

	if len(pm.queue) == 0 {
		return nil
	}

	pm.lock.Lock()
	firstMsg := pm.queue[0]
	firstAck := &proto.ScalarClockAck{
		ScalarClock: firstMsg.ScalarClock,
		From:        firstMsg.From,
	}

	if firstMsg.From == currUser {
		expectedAcks = len(pm.otherFroms)
	} else {
		expectedAcks = len(pm.otherFroms) - 1
	}

	log.Printf("Received %v/%v acks for [%v]\n", pm.receivedAcks[firstAck.String()], expectedAcks, firstMsg)

	if pm.receivedAcks[firstAck.String()] == expectedAcks && pm.thereAreMessagesFromAllInQueue(pm.queue[0].From) {
		deliverMsg = pm.queue[0]
		canDeliver = true
		pm.queue = pm.queue[1:]
		pm.presenceCounter[firstAck.From]--
	}
	pm.lock.Unlock()

	if canDeliver {
		return deliverMsg
	} else {
		return nil
	}
}

func (pm *PendingMessages) IncrementAckCounter(ack *proto.ScalarClockAck) {
	pm.lock.Lock()
	for i := range pm.queue {
		if isAckedMessage(pm.queue[i], ack) {
			pm.receivedAcks[ack.String()]++
		}
	}
	pm.lock.Unlock()
}
