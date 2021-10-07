package grpchandler

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

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

func (gh *GrpcHandler) SendMessages() error {
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
		newMessageBody := <-gh.peerStatus.RawMessageCh

		delay := rand.Intn(3000)
		log.Printf("Waiting %v millisec ...", delay)
		time.Sleep(time.Duration(delay) * time.Millisecond)

		_, err := c.SendToSequencer(context.Background(), &proto.UnorderedMessage{
			From: gh.peerStatus.CurrentUsername,
			Body: newMessageBody,
		})
		if err != nil {
			return err
		}
	}
}

func (s *ReceiverSequencerServer) SendToPeer(ctx context.Context, in *proto.OrderedMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (ID: %v)", in.GetBody(), in.GetFrom(), in.GetID())
	peer.InsertMessage(s.datastore, s.currentUser, in)
	return &empty.Empty{}, nil
}

func (gh *GrpcHandler) ReceiveMessages() error {
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
