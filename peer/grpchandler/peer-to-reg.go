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

func (gh *GrpcHandler) getOtherMembers(stream proto.Registration_SignClient) (bool, error) {
	for {
		member, err := stream.Recv()

		if err == io.EOF {
			return false, nil
		}

		if err != nil {
			errStatus, _ := status.FromError(err)
			if codes.InvalidArgument == errStatus.Code() {
				gh.peerStatus.InvalidCh <- true
				return true, nil
			} else {
				return false, err
			}
		}

		gh.peerStatus.Members = append(gh.peerStatus.Members, member)
	}
}

func (gh *GrpcHandler) SignToRegister() error {
	var currUser string

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", gh.registerAddr, gh.registerPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := proto.NewRegistrationClient(conn)

	for {
		currUser = <-gh.peerStatus.UsernameCh
		stream, err := c.Sign(context.Background(), &proto.ClientInfo{
			Username: currUser,
			Hostname: os.Getenv("HOSTNAME"),
		})
		if err != nil {
			return err
		}

		loopAgain, err := gh.getOtherMembers(stream)
		if err != nil {
			return nil
		}

		if !loopAgain {
			gh.peerStatus.CurrentUsername = currUser
			gh.peerStatus.DoneCh <- true
			return nil
		}
	}
}
