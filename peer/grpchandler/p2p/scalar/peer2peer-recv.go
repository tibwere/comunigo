package scalar

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

func (h *P2PScalarGRPCHandler) ReceiveMessages() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", h.comunicationPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	proto.RegisterComunigoServer(grpcServer, h)
	grpcServer.Serve(lis)

	return nil
}

func (h *P2PScalarGRPCHandler) tryToDeliverToDatastore() {
	if mess := h.pendingMsg.CheckIfIsReadyToDelivered(h.peerStatus.CurrentUsername); mess != nil {
		log.Printf("Delivered new message (Clock: %v - From: %v)\n", mess.GetScalarClock(), mess.GetFrom())
		peer.InsertScalarClockMessage(h.peerStatus.Datastore, h.peerStatus.CurrentUsername, mess)
	}
}

func (h *P2PScalarGRPCHandler) SendAckP2PScalar(ctx context.Context, in *proto.ScalarClockAck) (*empty.Empty, error) {
	log.Printf("Received ACK for %v (from %v)\n", in.GetScalarClock(), in.GetFrom())
	h.pendingMsg.IncrementAckCounter(in)

	h.tryToDeliverToDatastore()

	return &empty.Empty{}, nil
}

func (h *P2PScalarGRPCHandler) SendUpdateP2PScalar(ctx context.Context, in *proto.ScalarClockMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (ID: %v)", in.GetBody(), in.GetFrom(), in.GetScalarClock())
	h.pendingMsg.Insert(in)

	h.lockScalar.Lock()
	// L = max(t, L)
	if h.scalarClock < in.ScalarClock {
		h.scalarClock = in.ScalarClock
	}
	// L += 1
	h.scalarClock++

	// Invio del riscontro per il pacchetto ricevuto
	ack := &proto.ScalarClockAck{
		ScalarClock: in.GetScalarClock(),
		From:        in.GetFrom(),
	}
	h.lockScalar.Unlock()

	for _, ch := range h.scalarAcksChs {
		ch <- ack
	}

	h.tryToDeliverToDatastore()

	return &empty.Empty{}, nil
}
