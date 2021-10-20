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

func (h *P2PHandler) SendAckP2PScalar(ctx context.Context, in *proto.ScalarClockAck) (*empty.Empty, error) {
	log.Printf("Received ACK for %v (from %v)\n", in.GetTimestamp(), in.GetFrom())

	// Incremento del counter degli ack
	h.sData.InsertNewAck(in)

	return &empty.Empty{}, nil
}

func (h *P2PHandler) SendUpdateP2PScalar(ctx context.Context, in *proto.ScalarClockMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (timestamp: %v)", in.GetBody(), in.GetFrom(), in.GetTimestamp())

	// Aggiornamento del clock a seguito della ricezione e generazione l'ack
	ack := h.sData.UpdateClockAtRecv(in)

	// Riscontro del messaggio ricevuto a tutti gli altri membri (simulazione di auto-invio)
	h.sData.AckToAll(ack, true)

	// Inserimento del messaggio all'interno della coda di messaggi pendenti
	h.sData.InsertNewMessage(in)

	return &empty.Empty{}, nil
}

func (h *P2PHandler) SendUpdateP2PVectorial(ctx context.Context, in *proto.VectorialClockMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (timestamp: %v)", in.GetBody(), in.GetFrom(), in.GetTimestamp())
	h.vData.InsertNewMessage(in)
	return &empty.Empty{}, nil
}
