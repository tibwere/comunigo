package grpchandler

import (
	"context"
	"fmt"
	"io"
	"os"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getOtherMembers(stream proto.Registration_SignClient, invalidCh chan bool, usernameCh chan string) (bool, error) {
	for {
		member, err := stream.Recv()

		if err == io.EOF {
			return false, nil
		}

		if err != nil {
			errStatus, _ := status.FromError(err)
			if codes.InvalidArgument == errStatus.Code() {
				invalidCh <- true
				return true, nil
			} else {
				return false, err
			}
		}

		usernameCh <- member.GetUsername()
	}
}

func SignToRegister(addr string, port uint16, usernameCh chan string, invalidCh chan bool) (string, error) {
	var currUser string

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

	for {
		currUser := <-usernameCh
		stream, err := c.Sign(context.Background(), &proto.ClientInfo{
			Username: currUser,
			Hostname: os.Getenv("HOSTNAME"),
		})
		if err != nil {
			return "", err
		}

		loopAgain, err := getOtherMembers(stream, invalidCh, usernameCh)
		if err != nil {
			return "", nil
		}

		if !loopAgain {
			break
		}
	}

	return currUser, nil
}
