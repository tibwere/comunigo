package grpchandler

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

type ReceiverSequencerServer struct {
	proto.UnimplementedChatServer
	datastore   *redis.Client
	currentUser string
}

func (gh *GrpcHandler) SendMessages(wg *sync.WaitGroup) error {
	defer wg.Done()

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", gh.sequencerAddr, gh.sequencerPort),
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
			From: gh.peerStatus.CurrentUsername,
			Body: <-gh.peerStatus.RawMessageCh,
		})
		if err != nil {
			return err
		}
	}
}

func (s *ReceiverSequencerServer) SendToPeer(ctx context.Context, in *proto.OrderedMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (ID: %v)", in.GetBody(), in.GetFrom(), in.GetID())
	peer.InsertMessage(s.datastore, s.currentUser, in)
	fmt.Println("L'ho bello che inserito sto messaggio")
	return &empty.Empty{}, nil
}

func (gh *GrpcHandler) ReceiveMessages(wg *sync.WaitGroup) error {
	defer wg.Done()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", gh.sequencerPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()

	proto.RegisterChatServer(grpcServer, &ReceiverSequencerServer{
		datastore:   gh.peerStatus.Datastore,
		currentUser: gh.peerStatus.CurrentUsername,
	})

	grpcServer.Serve(lis)

	return nil
}
