package seq

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func (h *ToSequencerGRPCHandler) SendMessagesToSequencer(ctx context.Context) error {
	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", h.sequencerAddr, h.comunicationPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := proto.NewComunigoClient(conn)

	for {
		select {
		case <-ctx.Done():
			log.Println("Message sender to sequencer shutdown")
			return fmt.Errorf("signal caught")
		case newMessageBody := <-h.peerStatus.RawMessageCh:
			peer.WaitBeforeSend()
			_, err := c.SendFromPeerToSequencer(context.Background(), &proto.RawMessage{
				From: h.peerStatus.CurrentUsername,
				Body: newMessageBody,
			})
			if err != nil {
				return err
			}
		}

	}
}
