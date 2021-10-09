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

type SequencerServer struct {
	proto.UnimplementedComunigoServer
	sequenceNumber uint64
	seqCh          chan *proto.RawMessage
	connections    map[string]chan *proto.SequencerMessage
	port           uint16
	chatGroupSize  uint16
}

func (s *SequencerServer) LoadMembers(membersCh chan string, grpcServerToGetPeers *grpc.Server) {
	for i := 0; i < int(s.chatGroupSize); i++ {
		currentMember := <-membersCh
		s.connections[currentMember] = make(chan *proto.SequencerMessage)
		go s.sendBackMessages(currentMember)
	}

	grpcServerToGetPeers.GracefulStop()
}

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

func NewSequencerServer(port uint16, size uint16) *SequencerServer {

	seq := &SequencerServer{
		sequenceNumber: 0,
		seqCh:          make(chan *proto.RawMessage),
		connections:    make(map[string]chan *proto.SequencerMessage),
		port:           port,
		chatGroupSize:  size,
	}

	return seq
}

func ServePeers(seqServer *SequencerServer) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", seqServer.port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()

	proto.RegisterComunigoServer(grpcServer, seqServer)
	grpcServer.Serve(lis)

	return nil
}
