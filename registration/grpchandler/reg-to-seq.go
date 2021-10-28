package grpchandler

import (
	"context"
	"fmt"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

// Funzione che inizializza il sequencer inviandogli la lista dei peer connessi
// al gruppo di multicast tramite un'apposita procedura remota
func InitializeSequencer(addr string, port uint16, members []*proto.PeerInfo) error {

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", addr, port),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := proto.NewRegistrationClient(conn)

	stream, err := c.ExchangePeerInfoFromRegToSeq(context.Background())
	if err != nil {
		return err
	}

	for _, member := range members {
		if err := stream.Send(member); err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	return err
}
