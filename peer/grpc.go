package peer

import (
	"context"
	"fmt"
	"os"

	pb "gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func SignToRegister(addr string, port uint16, usernameCh chan string) error {

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", addr, port),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := pb.NewRegistrationClient(conn)

	res, err := c.Sign(context.Background(), &pb.ClientInfo{
		Username: <-usernameCh,
		Hostname: os.Getenv("HOSTNAME"),
	})
	if err != nil {
		return err
	}

	for _, client := range res.GetInfos() {
		usernameCh <- client.GetUsername()
	}

	return nil
}
