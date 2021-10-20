package scalar

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func (h *P2PScalarGRPCHandler) ReceiveMessages(ctx context.Context) error {
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

func (h *P2PScalarGRPCHandler) SendAckP2PScalar(ctx context.Context, in *proto.ScalarClockAck) (*empty.Empty, error) {
	log.Printf("Received ACK for %v (from %v)\n", in.GetTimestamp(), in.GetFrom())
	h.newAckCh <- in

	return &empty.Empty{}, nil
}

func (h *P2PScalarGRPCHandler) SendUpdateP2PScalar(ctx context.Context, in *proto.ScalarClockMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (timestamp: %v, current clock: %v)", in.GetBody(), in.GetFrom(), in.GetTimestamp(), h.clock)

	h.clockMu.Lock()
	// L = max(t, L)
	if h.clock < in.Timestamp {
		h.clock = in.Timestamp
	}

	// L += 1
	h.clock++

	// Invio del riscontro per il pacchetto ricevuto
	ack := &proto.ScalarClockAck{
		Timestamp: in.GetTimestamp(),
		From:      in.GetFrom(),
	}

	log.Printf("New clock value after update: %v\n", h.clock)
	h.clockMu.Unlock()

	// Incremento del contatore del mittente dell'ack perché non lo
	// invia a se stesso
	log.Printf("Autoincrement ACK counter after reception of message from %v\n", in.GetFrom())
	h.newAckCh <- ack

	for _, ch := range h.scalarAcksChs {
		ch <- ack
	}

	h.newMessageCh <- in

	return &empty.Empty{}, nil
}
