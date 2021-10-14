package vectorial

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func (h *P2PVectorialGRPCHandler) encapsulateMessage(body string) *proto.VectorialClockMessage {
	h.clockMu.Lock()
	h.incrementClockUnlocked(h.peerStatus.CurrentUsername)
	newMessage := &proto.VectorialClockMessage{
		Timestamp: h.vectorialClock,
		From:      h.peerStatus.CurrentUsername,
		Body:      body,
	}
	h.clockMu.Unlock()

	return newMessage
}

func (h *P2PVectorialGRPCHandler) MultiplexMessages() {

	for {
		newMessageBody := <-h.peerStatus.RawMessageCh

		log.Printf("Received from frontend: %v\n", newMessageBody)
		newMessage := h.encapsulateMessage(newMessageBody)
		log.Printf("Created new message with scalar clock %v\n", newMessage.GetTimestamp())

		for _, ch := range h.vectorialMessagesChs {
			ch <- newMessage
		}

		// Questo messaggio può essere direttamente consegnato perché di sicuro
		// rispetta la causalità
		peer.InsertVectorialClockMessage(h.peerStatus.Datastore, h.peerStatus.CurrentUsername, newMessage)
	}
}

func (h *P2PVectorialGRPCHandler) ConnectToPeers() error {
	errCh := make(chan error)

	for i := range h.peerStatus.OtherMembers {
		go func(index int, errCh chan error) {
			h.sendMessagesToOtherPeers(index, errCh)
		}(i, errCh)
	}

	return <-errCh
}

func (h *P2PVectorialGRPCHandler) sendMessagesToOtherPeers(index int, errCh chan error) {

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", h.peerStatus.OtherMembers[index].GetAddress(), h.comunicationPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		errCh <- err
		return
	}
	defer conn.Close()

	log.Printf("Succesfully linked to %v@%v\n", h.peerStatus.OtherMembers[index].GetUsername(), h.peerStatus.OtherMembers[index].GetAddress())

	c := proto.NewComunigoClient(conn)

	for {
		newMessage := <-h.vectorialMessagesChs[index]
		peer.WaitBeforeSend()
		log.Printf("Sending [%v] to %v@%v\n", newMessage, h.peerStatus.OtherMembers[index].Username, h.peerStatus.OtherMembers[index].Address)
		_, err := c.SendUpdateP2PVectorial(context.Background(), newMessage)
		if err != nil {
			errCh <- err
			return
		}
	}
}
