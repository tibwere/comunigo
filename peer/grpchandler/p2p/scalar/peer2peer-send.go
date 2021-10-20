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
	h.clockMu.Lock()
	h.clock++
	h.clockMu.Unlock()

	h.newMessageCh <- new
}

func (h *P2PScalarGRPCHandler) MultiplexMessages(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			log.Println("Messages multiplexer shutdown")
			return
		case newMessageBody := <-h.peerStatus.FrontBackCh:

			log.Printf("Received from frontend: %v\n", newMessageBody)
			h.clockMu.Lock()
			h.clock++
			newMessage := &proto.ScalarClockMessage{
				Timestamp: h.clock,
				From:      h.peerStatus.CurrentUsername,
				Body:      newMessageBody,
			}
			h.clockMu.Unlock()
			log.Printf("Created new message with scalar clock %v\n", newMessage.GetTimestamp())

			for _, ch := range h.scalarMessagesChs {
				ch <- newMessage
			}

			h.simulateRecv(newMessage)
		}
	}
}

func (h *P2PScalarGRPCHandler) ConnectToPeers(ctx context.Context) error {
	errCh := make(chan error)

	for i := range h.peerStatus.OtherMembers {
		index := i
		go func() {
			err := h.sendMessagesToOtherPeers(ctx, index)
			if err != nil {
				errCh <- err
			}
		}()
	}

	errMsg := ""
	for i := range h.peerStatus.OtherMembers {
		if len(errMsg) != 0 {
			errMsg = fmt.Sprintf("%v, %v->%v", errMsg, i, <-errCh)
		} else {
			errMsg = fmt.Sprintf("%v->%v", i, <-errCh)
		}
	}

	return fmt.Errorf(errMsg)
}

func (h *P2PScalarGRPCHandler) sendMessagesToOtherPeers(ctx context.Context, index int) error {

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", h.peerStatus.OtherMembers[index].GetAddress(), h.comunicationPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("Succesfully linked to %v@%v\n", h.peerStatus.OtherMembers[index].GetUsername(), h.peerStatus.OtherMembers[index].GetAddress())

	c := proto.NewComunigoClient(conn)

	for {
		var newMessage *proto.ScalarClockMessage
		var newAck *proto.ScalarClockAck

		select {
		case <-ctx.Done():
			log.Printf("Message sender %v shutdown\n", index)
			return fmt.Errorf("signal caught")

		case newMessage = <-h.scalarMessagesChs[index]:
			peer.WaitBeforeSend()
			log.Printf("Sending [%v] to %v@%v\n", newMessage, h.peerStatus.OtherMembers[index].Username, h.peerStatus.OtherMembers[index].Address)
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
