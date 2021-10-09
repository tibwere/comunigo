package seq

import (
	"context"
	"fmt"

	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func (h *ToSequencerGRPCHandler) SendMessagesToSequencer() error {
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
		newMessageBody := <-h.peerStatus.RawMessageCh

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
