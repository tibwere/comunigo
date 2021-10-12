package grpchandler

import (
	"fmt"
	"io"
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

func (s *FromRegisterServer) GetPeersFromRegister(port uint16, fromRegToSeqGRPCserver *grpc.Server) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}

	proto.RegisterRegistrationServer(fromRegToSeqGRPCserver, s)
	fromRegToSeqGRPCserver.Serve(lis)

	return nil
}
