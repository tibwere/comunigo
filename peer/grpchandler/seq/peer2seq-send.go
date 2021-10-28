package seq

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

// "Metodo della classe ToSequencerGRPCHandler" che permette l'invio
// dei messaggi al sequencer
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
		case newMessageBody := <-h.peerStatus.GetFromFrontendBackendChannel():
			_, err := c.SendFromPeerToSequencer(context.Background(), &proto.RawMessage{
				From: h.peerStatus.GetCurrentUsername(),
				Body: newMessageBody,
			})
			if err != nil {
				return err
			}
		}

	}
}
