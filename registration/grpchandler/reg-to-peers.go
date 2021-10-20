package grpchandler

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type exchangeInformationFromUpdaterAndHandler struct {
	clientInfo        *proto.PeerInfo
	isUsernameValidCh chan bool
}

type RegistrationServer struct {
	proto.UnimplementedRegistrationServer
	newMemberCh       chan *exchangeInformationFromUpdaterAndHandler
	memberInformation []*exchangeInformationFromUpdaterAndHandler
	numberOfClients   uint16
}

func NewRegistrationServer(size uint16) *RegistrationServer {
	return &RegistrationServer{
		newMemberCh:       make(chan *exchangeInformationFromUpdaterAndHandler),
		memberInformation: []*exchangeInformationFromUpdaterAndHandler{},
		numberOfClients:   size,
	}
}

func (s *RegistrationServer) getChatGroupMembers() []*proto.PeerInfo {
	var members []*proto.PeerInfo

	for _, info := range s.memberInformation {
		members = append(members, info.clientInfo)
	}

	return members
}

func (s *RegistrationServer) getValidityList() []chan bool {
	var validityList []chan bool

	for _, info := range s.memberInformation {
		validityList = append(validityList, info.isUsernameValidCh)
	}

	return validityList
}

func (s *RegistrationServer) isValidUsername(username string) bool {
	for _, member := range s.memberInformation {
		if member.clientInfo.GetUsername() == username {
			return false
		}
	}
	return true
}

func (s *RegistrationServer) UpdateMembers(ctx context.Context, grpcServer *grpc.Server, seqAddr string, seqPort uint16, needSequencer bool) error {
	// stop del server al completamento del gruppo di multicast
	defer grpcServer.GracefulStop()

	for uint16(len(s.memberInformation)) < s.numberOfClients {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case info := <-s.newMemberCh:
			if s.isValidUsername(info.clientInfo.GetUsername()) {
				s.memberInformation = append(s.memberInformation, info)
			} else {
				info.isUsernameValidCh <- false
			}
		}
	}

	if needSequencer {
		InitializeSequencer(seqAddr, seqPort, s.getChatGroupMembers())
	}

	// Invia un messaggio di sincronizzazione a tutte
	// le goroutine per inviare la risposta al client
	for _, ch := range s.getValidityList() {
		ch <- true
	}

	return nil
}

func (s *RegistrationServer) Sign(in *proto.NewUser, stream proto.Registration_SignServer) error {

	p, _ := peer.FromContext(stream.Context())
	peerIP := strings.Split(p.Addr.String(), ":")[0]
	log.Printf("Received request from: %v@%v\n", in.GetUsername(), peerIP)

	// Crea un canale con cui comprendere se l'username è valido o meno
	validityCh := make(chan bool)

	// Invia i metadati relativi al nuovo utente alla goroutine dedicata all'aggiornamento
	// del datastore
	s.newMemberCh <- &exchangeInformationFromUpdaterAndHandler{
		clientInfo: &proto.PeerInfo{
			Username: in.GetUsername(),
			Address:  peerIP,
		},
		isUsernameValidCh: validityCh,
	}

	// Aspetta il messaggio di sincronizzazione per restituire le info al chiamante
	isValid := <-validityCh

	if isValid {
		for _, member := range s.getChatGroupMembers() {
			if err := stream.Send(member); err != nil {
				return err
			}
		}
		return nil
	} else {
		return status.Errorf(codes.InvalidArgument, "Username already in use, please retry with another one!")
	}
}

func ServeSignRequests(ctx context.Context, exposedPort uint16, regServer *RegistrationServer, grpcServer *grpc.Server) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", exposedPort))
	if err != nil {
		return err
	}

	proto.RegisterRegistrationServer(grpcServer, regServer)
	go grpcServer.Serve(lis)

	<-ctx.Done()
	log.Println("Registration server shutdown")
	grpcServer.GracefulStop()
	return fmt.Errorf("signal caught")
}
