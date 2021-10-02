package grpchandler

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

type StartupSequencerServer struct {
	proto.UnimplementedRegistrationServer
	MembersCh chan string
}

func (s *StartupSequencerServer) StartSequencer(stream proto.Registration_StartSequencerServer) error {

	for {
		member, err := stream.Recv()
		if err == io.EOF {

			return stream.SendAndClose(&empty.Empty{})
		}
		if err != nil {
			return err
		}

		s.MembersCh <- member.GetHostname()
	}
}

func GetClientsFromRegister(port uint16, startupServer *StartupSequencerServer, grpcServer *grpc.Server, wg *sync.WaitGroup) {
	defer wg.Done()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}

	proto.RegisterRegistrationServer(grpcServer, startupServer)
	grpcServer.Serve(lis)
}
