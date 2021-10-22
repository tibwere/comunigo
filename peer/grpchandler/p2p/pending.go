package p2p

import (
	"context"
	"fmt"
	"log"
)

func (h *P2PHandler) MessageQueueHandler(ctx context.Context) error {
	if h.modality == P2P_SCALAR {
		return h.messageQueueHandlerSC(ctx)
	} else {
		return h.messageQueueHandlerVC(ctx)
	}
}

func (h *P2PHandler) messageQueueHandlerSC(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case newMessage := <-h.sData.GetMessageCh():
			log.Printf("Insert [%v] into pendant queue\n", newMessage)
			h.sData.PushIntoPendingList(newMessage)
			h.sData.SyncDatastore(h.peerStatus.Datastore, h.peerStatus.CurrentUsername, h.peerStatus.OtherMembers)
		case newAck := <-h.sData.GetAckCh():
			log.Printf("Increment ack counter of [%v:%v]\n", newAck.From, newAck.Timestamp)
			h.sData.IncrementAckCounter(newAck)
			h.sData.SyncDatastore(h.peerStatus.Datastore, h.peerStatus.CurrentUsername, h.peerStatus.OtherMembers)
		}
	}
}

func (h *P2PHandler) messageQueueHandlerVC(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case newMessage := <-h.vData.GetReceivedCh():
			log.Printf("Insert [%v] into pendant queue\n", newMessage)
			h.vData.PushIntoPendingList(newMessage)
			if err := h.vData.SyncDatastore(h.peerStatus.Datastore, h.peerStatus.CurrentUsername); err != nil {
				return err
			}
		}

	}
}
