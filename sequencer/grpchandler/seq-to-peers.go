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

type SequencerServer struct {
	proto.UnimplementedChatServer
	sequenceNumber uint64
	seqCh          chan *proto.UnorderedMessage
	members        []string
	port           uint16
	chatGroupSize  uint16
}

func (s *SequencerServer) LoadMembers(membersCh chan string, grpcServer *grpc.Server, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < int(s.chatGroupSize); i++ {
		s.members = append(s.members, <-membersCh)
	}

	grpcServer.GracefulStop()
}

func (s *SequencerServer) SendToSequencer(ctx context.Context, in *proto.UnorderedMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v\n", in.GetBody(), in.GetFrom())
	s.seqCh <- in
	return &empty.Empty{}, nil
}

func (s *SequencerServer) OrderMessages(port uint16) {
	for {
		unordered := <-s.seqCh

		ordered := &proto.OrderedMessage{
			ID:   s.sequenceNumber,
			From: unordered.GetFrom(),
			Body: unordered.GetBody(),
		}
		s.sequenceNumber++

		s.sendBackToPeers(ordered, port)
	}
}

func NewSequencerserver(port uint16, size uint16) *SequencerServer {

	seq := &SequencerServer{
		sequenceNumber: 0,
		seqCh:          make(chan *proto.UnorderedMessage),
		members:        []string{},
		port:           port,
		chatGroupSize:  size,
	}

	return seq
}

func (s *SequencerServer) sendBackToPeers(ordered *proto.OrderedMessage, port uint16) error {

	for _, peer := range s.members {
		conn, err := grpc.Dial(
			fmt.Sprintf("%v:%v", peer, port),
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)
		if err != nil {
			return err
		}
		defer conn.Close()

		c := proto.NewChatClient(conn)

		_, err = c.SendToPeer(context.Background(), ordered)
		if err != nil {
			return err
		}
	}

	return nil
}

func ServePeers(seqServer *SequencerServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", seqServer.port))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()

	proto.RegisterChatServer(grpcServer, seqServer)
	grpcServer.Serve(lis)
}
