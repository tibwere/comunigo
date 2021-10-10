package grpchandler

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/proto"
	"golang.org/x/sync/errgroup"
)

func (s *SequencerServer) ExchangePeerInfoFromRegToSeq(stream proto.Registration_ExchangePeerInfoFromRegToSeqServer) error {
	errs, _ := errgroup.WithContext(context.Background())
	for {
		member, err := stream.Recv()
		if err == io.EOF {
			if err := stream.SendAndClose(&empty.Empty{}); err != nil {
				return err
			}
			break
		}
		if err != nil {
			return err
		}

		s.connections[member.Address] = make(chan *proto.SequencerMessage)
		errs.Go(func() error {
			return s.sendBackMessages(member.Address)
		})
	}

	s.fromRegToSeqGRPCserver.GracefulStop()
	return errs.Wait()
}

func (s *SequencerServer) GetPeersFromRegister(port uint16) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}

	proto.RegisterRegistrationServer(s.fromRegToSeqGRPCserver, s)
	s.fromRegToSeqGRPCserver.Serve(lis)
	return nil
}
