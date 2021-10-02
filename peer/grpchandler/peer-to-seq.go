package grpchandler

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

type ReceiverSequencerServer struct {
	proto.UnimplementedChatServer
}

func SendMessages(addr string, port uint16, currUser string, messageCh chan string, wg *sync.WaitGroup) error {
	defer wg.Done()

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", addr, port),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := proto.NewChatClient(conn)

	for {
		_, err := c.SendToSequencer(context.Background(), &proto.UnorderedMessage{
			From: currUser,
			Body: <-messageCh,
		})
		if err != nil {
			return err
		}
	}
}

func (s *ReceiverSequencerServer) SendToPeer(ctx context.Context, in *proto.OrderedMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (ID: %v)", in.GetBody(), in.GetFrom(), in.GetID())
	return &empty.Empty{}, nil
}

func ReceiveMessages(port uint16, wg *sync.WaitGroup) error {
	defer wg.Done()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()

	proto.RegisterChatServer(grpcServer, &ReceiverSequencerServer{})
	grpcServer.Serve(lis)

	return nil
}
