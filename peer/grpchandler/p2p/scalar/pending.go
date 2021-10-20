package scalar

import (
	"context"
	"fmt"
	"log"
	"sort"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
)

func (h *P2PScalarGRPCHandler) MessageQueueHandler(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case newMessage := <-h.newMessageCh:
			log.Printf("Insert [%v] into pendant queue\n", newMessage)
			h.pendingMsg = append(h.pendingMsg, newMessage)

			sort.Slice(h.pendingMsg, func(i, j int) bool {

				iClock := h.pendingMsg[i].Timestamp
				jClock := h.pendingMsg[j].Timestamp
				iFrom := h.pendingMsg[i].From
				jFrom := h.pendingMsg[i].From

				return iClock < jClock || (iClock == jClock && iFrom < jFrom)
			})
			h.presenceCounter[newMessage.From]++
			h.syncDatastore()

		case newAck := <-h.newAckCh:
			h.receivedAcks[newAck.String()]++
			h.syncDatastore()
		}
	}
}

func (h *P2PScalarGRPCHandler) syncDatastore() {
	for _, mess := range h.deliverMessagesIfPossible() {
		log.Printf("Delivered new message (Clock: %v - From: %v)\n", mess.GetTimestamp(), mess.GetFrom())
		peer.InsertScalarClockMessage(h.peerStatus.Datastore, h.peerStatus.CurrentUsername, mess)
	}
}

func (h *P2PScalarGRPCHandler) deliverMessagesIfPossible() []*proto.ScalarClockMessage {
	var deliverList []*proto.ScalarClockMessage

	if len(h.pendingMsg) == 0 {
		return deliverList
	}

	firstMsg := h.pendingMsg[0]
	firstAck := &proto.ScalarClockAck{
		Timestamp: firstMsg.Timestamp,
		From:      firstMsg.From,
	}
	nMember := len(h.peerStatus.OtherMembers)

	log.Printf("Received %v/%v acks for [%v]\n", h.receivedAcks[firstAck.String()], nMember, firstMsg)

	if h.receivedAcks[firstAck.String()] == nMember && h.thereAreMessagesFromAllInQueue(firstMsg.From) {
		deliverList = append(deliverList, firstMsg)
		h.pendingMsg = h.pendingMsg[1:]
		h.presenceCounter[firstAck.From]--
	}

	return deliverList
}

func (h *P2PScalarGRPCHandler) thereAreMessagesFromAllInQueue(actualFrom string) bool {
	for from, presences := range h.presenceCounter {
		if from != actualFrom && presences == 0 {
			log.Printf("Member %v does not have messages in queue at this moment. Cannot deliver message\n", from)
			return false
		}
	}

	log.Println("All other members has at least a message in queue at this moment.")

	return true
}
