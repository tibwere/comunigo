package grpchandler

import (
	"context"
	"fmt"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func InitializeSequencer(addr string, port uint16, members []*proto.ClientInfo) error {

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

	stream, err := c.StartSequencer(context.Background())
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
