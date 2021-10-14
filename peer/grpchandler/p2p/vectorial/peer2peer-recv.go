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

func (h *P2PVectorialGRPCHandler) MessageQueueHandler() error {
	for {
		newMessage := <-h.receivedCh

		h.pendingMsg = append(h.pendingMsg, newMessage)
		deliverables := h.tryToDeliverToDatastore()
		log.Printf("%v new message can be delivered\n", len(deliverables))

		if len(deliverables) != 0 {
			for _, mess := range deliverables {
				log.Printf("Delivering new message: Timestamp: %v, From: %v, Body: '%v'\n", mess.Timestamp, mess.From, mess.Body)

				h.clockMu.Lock()
				h.incrementClockUnlocked(mess.From)
				h.clockMu.Unlock()

				if err := peer.InsertVectorialClockMessage(h.peerStatus.Datastore, h.peerStatus.CurrentUsername, mess); err != nil {
					return err
				}
			}
		}

	}
}

func (h *P2PVectorialGRPCHandler) SendUpdateP2PVectorial(ctx context.Context, in *proto.VectorialClockMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v (timestamp: %v, current clock: %v)", in.GetBody(), in.GetFrom(), in.GetTimestamp(), h.vectorialClock)
	h.receivedCh <- in
	return &empty.Empty{}, nil
}

func (h *P2PVectorialGRPCHandler) tryToDeliverToDatastore() []*proto.VectorialClockMessage {

	var deliverables []*proto.VectorialClockMessage

	for _, mess := range h.pendingMsg {
		if h.isDeliverable(mess) {
			deliverables = append(deliverables, mess)
		}
	}

	return deliverables
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
