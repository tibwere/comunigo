package grpchandler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func (s *FromRegisterServer) ExchangePeerInfoFromRegToSeq(stream proto.Registration_ExchangePeerInfoFromRegToSeqServer) error {
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

		s.memberCh <- member
	}

	return nil
}

func (s *FromRegisterServer) GetPeersFromRegister(ctx context.Context, port uint16, fromRegToSeqGRPCserver *grpc.Server) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}

	proto.RegisterRegistrationServer(fromRegToSeqGRPCserver, s)
	go fromRegToSeqGRPCserver.Serve(lis)

	<-ctx.Done()
	log.Println("Sequencer server shutdown")
	fromRegToSeqGRPCserver.GracefulStop()
	return fmt.Errorf("signal caught")
}
