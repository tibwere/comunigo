package grpchandler

import (
	"context"
	"fmt"
	"io"
	"os"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func SignToRegister(addr string, port uint16, usernameCh chan string) (string, error) {

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", addr, port),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := proto.NewRegistrationClient(conn)

	currUser := <-usernameCh
	stream, err := c.Sign(context.Background(), &proto.ClientInfo{
		Username: currUser,
		Hostname: os.Getenv("HOSTNAME"),
	})
	if err != nil {
		return "", err
	}

	for {
		member, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		usernameCh <- member.GetUsername()
	}

	return currUser, nil
}
