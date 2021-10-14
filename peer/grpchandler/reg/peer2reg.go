package reg

import (
	"context"
	"fmt"
	"io"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *ToRegisterGRPCHandler) getOtherMembers(currUser string, stream proto.Registration_SignClient) (bool, error) {
	for {
		member, err := stream.Recv()

		if err == io.EOF {
			return false, nil
		}

		if err != nil {
			errStatus, _ := status.FromError(err)
			if codes.InvalidArgument == errStatus.Code() {
				h.peerStatus.InvalidCh <- true
				return true, nil
			} else {
				return false, err
			}
		}

		if currUser != member.GetUsername() {
			h.peerStatus.OtherMembers = append(h.peerStatus.OtherMembers, member)
		}
	}
}

func (h *ToRegisterGRPCHandler) SignToRegister() error {
	var currUser string

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", h.registerAddr, h.registerPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := proto.NewRegistrationClient(conn)

	for {
		currUser = <-h.peerStatus.UsernameCh
		stream, err := c.Sign(context.Background(), &proto.NewUser{
			Username: currUser,
		})
		if err != nil {
			return err
		}

		loopAgain, err := h.getOtherMembers(currUser, stream)
		if err != nil {
			return nil
		}

		if !loopAgain {
			h.peerStatus.CurrentUsername = currUser
			h.peerStatus.DoneCh <- true
			return nil
		}
	}
}
