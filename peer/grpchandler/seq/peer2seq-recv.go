package seq

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

// "Metodo della classe ToSequencerGRPCHandler" per l'implementazione della RPC SendFromSequencerToPeer server-side
func (s *ToSequencerGRPCHandler) SendFromSequencerToPeer(ctx context.Context, in *proto.SequencerMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (ID: %v)", in.GetBody(), in.GetFrom(), in.GetTimestamp())
	s.peerStatus.RPUSHMessage(in)
	return &empty.Empty{}, nil
}

// "Metodo della classe ToSequencerGRPCHandler" che inizializza il server gRPC per la ricezione
// dei messaggi provenienti dal sequencer
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
