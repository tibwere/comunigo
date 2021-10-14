package vectorial

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

func (h *P2PVectorialGRPCHandler) ReceiveMessages() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", h.comunicationPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	proto.RegisterComunigoServer(grpcServer, h)
	grpcServer.Serve(lis)

	return nil
}

func (h *P2PVectorialGRPCHandler) SendUpdateP2PVectorial(ctx context.Context, in *proto.VectorialClockMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (timestamp: %v, current clock: %v)", in.GetBody(), in.GetFrom(), in.GetTimestamp(), h.vectorialClock)
	h.insertNewMessage(in)
	deliverables := h.tryToDeliverToDatastore()
	log.Printf("New message can be delivered: %v\n", deliverables)

	if len(deliverables) != 0 {
		for _, mess := range deliverables {
			h.incrementClock(mess.From)
			peer.InsertVectorialClockMessage(h.peerStatus.Datastore, h.peerStatus.CurrentUsername, mess)
		}
	}

	return &empty.Empty{}, nil
}

func (h *P2PVectorialGRPCHandler) tryToDeliverToDatastore() []*proto.VectorialClockMessage {

	var deliverables []*proto.VectorialClockMessage

	h.mu.Lock()
	for _, mess := range h.pendingMsg {
		if h.isDeliverable(mess) {
			deliverables = append(deliverables, mess)
		}
	}
	h.mu.Unlock()

	return deliverables
}

func (h *P2PVectorialGRPCHandler) insertNewMessage(mess *proto.VectorialClockMessage) {
	h.mu.Lock()
	h.pendingMsg = append(h.pendingMsg, mess)
	h.mu.Unlock()
}

func (h *P2PVectorialGRPCHandler) isDeliverable(mess *proto.VectorialClockMessage) bool {

	fromIndex := h.memberIndexs[mess.From]
	for i, messEntry := range mess.Timestamp {
		localEntry := h.vectorialClock[i]
		if i == fromIndex && messEntry != localEntry+1 {
			return false
		}

		if i != fromIndex && messEntry > localEntry {
			return false
		}
	}

	return true
}
