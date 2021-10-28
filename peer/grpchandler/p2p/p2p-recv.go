package p2p

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

// "Metodo della classe P2PHandler" che inizializza il server gRPC per la ricezione
// dei messaggi dagli altri peer partecipanti al gruppo di multicast
func (h *P2PHandler) ReceiveMessages(ctx context.Context) error {
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

// "Metodo della classe P2PHandler" per l'implementazione della RPC SendAckP2PScalar server-side
func (h *P2PHandler) SendAckP2PScalar(ctx context.Context, in *proto.ScalarClockAck) (*empty.Empty, error) {
	log.Printf("Received ACK for %v (from %v)\n", in.GetTimestamp(), in.GetFrom())
	h.sData.InsertNewAck(in)
	return &empty.Empty{}, nil
}

// "Metodo della classe P2PHandler" per l'implementazione della RPC SendUpdateP2PScalar server-side
func (h *P2PHandler) SendUpdateP2PScalar(ctx context.Context, in *proto.ScalarClockMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (timestamp: %v)", in.GetBody(), in.GetFrom(), in.GetTimestamp())
	ack := h.sData.UpdateClockAtRecv(in)
	h.sData.AckToAll(ack)
	h.sData.InsertNewMessage(in)
	return &empty.Empty{}, nil
}

// "Metodo della classe P2PHandler" per l'implementazione della RPC SendUpdateP2PVectorial server-side
func (h *P2PHandler) SendUpdateP2PVectorial(ctx context.Context, in *proto.VectorialClockMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (timestamp: %v)", in.GetBody(), in.GetFrom(), in.GetTimestamp())
	h.vData.InsertNewMessage(in)
	return &empty.Empty{}, nil
}
