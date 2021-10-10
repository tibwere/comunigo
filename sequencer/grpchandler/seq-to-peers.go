package grpchandler

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func (s *SequencerServer) sendBackMessages(addr string) error {
	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", addr, s.port),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := proto.NewComunigoClient(conn)

	for {
		_, err = c.SendFromSequencerToPeer(context.Background(), <-s.connections[addr])
		if err != nil {
			return err
		}
	}
}

func (s *SequencerServer) SendFromPeerToSequencer(ctx context.Context, in *proto.RawMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v\n", in.GetBody(), in.GetFrom())
	s.seqCh <- in
	return &empty.Empty{}, nil
}

func (s *SequencerServer) OrderMessages() {
	for {
		unordered := <-s.seqCh

		ordered := &proto.SequencerMessage{
			SequenceNumber: s.sequenceNumber,
			From:           unordered.GetFrom(),
			Body:           unordered.GetBody(),
		}
		s.sequenceNumber++

		for _, ch := range s.connections {
			ch <- ordered
		}
	}
}

func (s *SequencerServer) ServePeers() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	proto.RegisterComunigoServer(grpcServer, s)
	grpcServer.Serve(lis)
	return nil
}
