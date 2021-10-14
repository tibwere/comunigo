package seq

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func (s *ToSequencerGRPCHandler) SendFromSequencerToPeer(ctx context.Context, in *proto.SequencerMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (ID: %v)", in.GetBody(), in.GetFrom(), in.GetTimestamp())
	peer.InsertSequencerMessage(s.peerStatus.Datastore, s.peerStatus.CurrentUsername, in)
	return &empty.Empty{}, nil
}

func (h *ToSequencerGRPCHandler) ReceiveMessages() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", h.comunicationPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()

	proto.RegisterComunigoServer(grpcServer, h)

	grpcServer.Serve(lis)

	return nil
}
