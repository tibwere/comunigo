package scalar

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func (h *P2PScalarGRPCHandler) MultiplexMessages() {

	for {
		newMessageBody := <-h.peerStatus.RawMessageCh

		log.Printf("Received from frontend: %v\n", newMessageBody)
		h.lockScalar.Lock()
		h.scalarClock++
		newMessage := &proto.ScalarClockMessage{
			ScalarClock: h.scalarClock,
			From:        h.peerStatus.CurrentUsername,
			Body:        newMessageBody,
		}
		h.lockScalar.Unlock()
		log.Printf("Created new message with scalar clock %v\n", newMessage.GetScalarClock())

		h.pendingMsg.Insert(newMessage)

		for _, ch := range h.scalarMessagesChs {
			ch <- newMessage
		}
	}
}

func (h *P2PScalarGRPCHandler) ConnectToPeers() {
	for i := range h.peerStatus.Members {
		go h.sendMessagesToOtherPeers(i)
	}
}

func (h *P2PScalarGRPCHandler) sendMessagesToOtherPeers(index int) error {

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", h.peerStatus.Members[index].GetAddress(), h.comunicationPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("Succesfully linked to %v@%v\n", h.peerStatus.Members[index].GetUsername(), h.peerStatus.Members[index].GetAddress())

	c := proto.NewComunigoClient(conn)

	for {
		var newMessage *proto.ScalarClockMessage
		var newAck *proto.ScalarClockAck

		select {
		case newMessage = <-h.scalarMessagesChs[index]:
			peer.WaitBeforeSend()
			log.Printf("Sending [%v] to %v@%v\n", newMessage, h.peerStatus.Members[index].Username, h.peerStatus.Members[index].Address)
			_, err := c.SendUpdateP2PScalar(context.Background(), newMessage)
			if err != nil {
				return err
			}

		case newAck = <-h.scalarAcksChs[index]:
			peer.WaitBeforeSend()
			_, err := c.SendAckP2PScalar(context.Background(), newAck)
			if err != nil {
				return err
			}
		}
	}
}
