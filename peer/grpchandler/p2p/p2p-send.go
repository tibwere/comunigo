package p2p

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func (h *P2PHandler) MultiplexMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case newMessageBody := <-h.peerStatus.FrontBackCh:
			log.Printf("Received from frontend: %v\n", newMessageBody)

			if h.modality == P2P_SCALAR {
				newMessage := h.sData.GenerateNewMessage(h.peerStatus.CurrentUsername, newMessageBody)
				log.Printf("Created new message with scalar clock %v\n", newMessage.GetTimestamp())

				h.sData.SendToAll(newMessage)
			} else {
				newMessage := h.vData.GenerateNewMessage(h.peerStatus.CurrentUsername, newMessageBody)
				log.Printf("Created new message with vectorial clock %v\n", newMessage.GetTimestamp())

				if err := h.vData.SendToAll(newMessage, h.peerStatus.Datastore, h.peerStatus.CurrentUsername); err != nil {
					return err
				}
			}

		}
	}
}

func (h *P2PHandler) ConnectToPeers(ctx context.Context) error {
	errCh := make(chan error)

	for i := range h.peerStatus.OtherMembers {
		index := i
		go func() {
			err := h.sendToOther(ctx, index)
			if err != nil {
				errCh <- err
			}
		}()
	}

	errMsg := ""
	for _, m := range h.peerStatus.OtherMembers {
		errMsg += fmt.Sprintf("Handler for: %v->%v, ", m.Username, <-errCh)
	}
	// rimuove l'ulitmo ", "
	errMsg = errMsg[:len(errMsg)-2]

	return fmt.Errorf(errMsg)
}

func (h *P2PHandler) sendToOther(ctx context.Context, index int) error {

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

	if h.modality == P2P_SCALAR {
		return h.sendLoopSC(ctx, c, index)
	} else {
		return h.sendLoopVC(ctx, c, index)
	}
}

func (h *P2PHandler) sendLoopSC(ctx context.Context, c proto.ComunigoClient, index int) error {
	for {
		var newMessage *proto.ScalarClockMessage
		var newAck *proto.ScalarClockAck

		select {
		case <-ctx.Done():
			log.Printf("Message sender %v shutdown\n", index)
			return fmt.Errorf("signal caught")

		case newMessage = <-h.sData.GetIncomingMsgToBeSentCh(index):
			peer.WaitBeforeSend()
			log.Printf("Sending [%v] to %v@%v\n", newMessage, h.peerStatus.OtherMembers[index].Username, h.peerStatus.OtherMembers[index].Address)
			_, err := c.SendUpdateP2PScalar(context.Background(), newMessage)
			if err != nil {
				return err
			}

		case newAck = <-h.sData.GetIncomingAckToBeSentCh(index):
			peer.WaitBeforeSend()
			_, err := c.SendAckP2PScalar(context.Background(), newAck)
			if err != nil {
				return err
			}
		}
	}
}

func (h *P2PHandler) sendLoopVC(ctx context.Context, c proto.ComunigoClient, index int) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Message sender %v shutdown\n", index)
			return fmt.Errorf("signal caught")
		case newMessage := <-h.vData.GetIncomingMsgToBeSentCh(index):
			peer.WaitBeforeSend()
			log.Printf("Sending [%v] to %v@%v\n", newMessage, h.peerStatus.OtherMembers[index].Username, h.peerStatus.OtherMembers[index].Address)
			_, err := c.SendUpdateP2PVectorial(context.Background(), newMessage)
			if err != nil {
				return err
			}
		}
	}
}
