package scalar

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func (h *P2PScalarGRPCHandler) simulateRecv(new *proto.ScalarClockMessage) {
	h.lockScalar.Lock()
	h.scalarClock++
	h.lockScalar.Unlock()

	h.pendingMsg.Insert(new)
}

func (h *P2PScalarGRPCHandler) MultiplexMessages() {

	for {
		newMessageBody := <-h.peerStatus.RawMessageCh

		log.Printf("Received from frontend: %v\n", newMessageBody)
		h.lockScalar.Lock()
		h.scalarClock++
		newMessage := &proto.ScalarClockMessage{
			Timestamp: h.scalarClock,
			From:      h.peerStatus.CurrentUsername,
			Body:      newMessageBody,
		}
		h.lockScalar.Unlock()
		log.Printf("Created new message with scalar clock %v\n", newMessage.GetTimestamp())

		for _, ch := range h.scalarMessagesChs {
			ch <- newMessage
		}

		h.simulateRecv(newMessage)
	}
}

func (h *P2PScalarGRPCHandler) ConnectToPeers() error {
	errCh := make(chan error)

	for i := range h.peerStatus.OtherMembers {
		go func(index int, errCh chan error) {
			h.sendMessagesToOtherPeers(index, errCh)
		}(i, errCh)
	}

	return <-errCh
}

func (h *P2PScalarGRPCHandler) sendMessagesToOtherPeers(index int, errCh chan error) {

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
		var newMessage *proto.ScalarClockMessage
		var newAck *proto.ScalarClockAck

		select {
		case newMessage = <-h.scalarMessagesChs[index]:
			peer.WaitBeforeSend()
			log.Printf("Sending [%v] to %v@%v\n", newMessage, h.peerStatus.OtherMembers[index].Username, h.peerStatus.OtherMembers[index].Address)
			_, err := c.SendUpdateP2PScalar(context.Background(), newMessage)
			if err != nil {
				errCh <- err
				return
			}

		case newAck = <-h.scalarAcksChs[index]:
			peer.WaitBeforeSend()
			_, err := c.SendAckP2PScalar(context.Background(), newAck)
			if err != nil {
				errCh <- err
				return
			}
		}
	}
}
