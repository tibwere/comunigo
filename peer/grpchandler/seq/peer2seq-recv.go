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
	peer.RPUSHMessage(s.peerStatus.Datastore, s.peerStatus.CurrentUsername, in)
	return &empty.Empty{}, nil
}

func (h *ToSequencerGRPCHandler) ReceiveMessages(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", h.comunicationPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	proto.RegisterComunigoServer(grpcServer, h)

	go grpcServer.Serve(lis)

	<-ctx.Done()
	log.Println("Message receiver from sequencer shutdown")
	grpcServer.GracefulStop()
	return fmt.Errorf("signal caught")
}
