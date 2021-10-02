package grpchandler

import (
	"fmt"
	"log"
	"net"
	"sync"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type exchangeInformationFromUpdaterAndHandler struct {
	clientInfo        *proto.ClientInfo
	isUsernameValidCh chan bool
}

type RegistrationServer struct {
	proto.UnimplementedRegistrationServer
	newMemberCh       chan *exchangeInformationFromUpdaterAndHandler
	memberInformation []*exchangeInformationFromUpdaterAndHandler
	numberOfClients   uint16
}

func (s *RegistrationServer) getChatGroupMembers() []*proto.ClientInfo {
	var members []*proto.ClientInfo

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

func (s *RegistrationServer) UpdateMembers(grpcServer *grpc.Server, seqAddr string, seqPort uint16, wg *sync.WaitGroup) {
	defer wg.Done()

	for uint16(len(s.memberInformation)) < s.numberOfClients {
		info := <-s.newMemberCh
		if s.isValidUsername(info.clientInfo.GetUsername()) {
			s.memberInformation = append(s.memberInformation, info)
		} else {
			info.isUsernameValidCh <- false
		}
	}

	InitializeSequencer(seqAddr, seqPort, s.getChatGroupMembers())

	// Invia un messaggio di sincronizzazione a tutte
	// le goroutine per inviare la risposta al client
	for _, ch := range s.getValidityList() {
		ch <- true
	}

	// stop del server poiché è stato completato il gruppo di multicast
	grpcServer.GracefulStop()
}

func (s *RegistrationServer) Sign(in *proto.ClientInfo, stream proto.Registration_SignServer) error {
	log.Printf("Received request from: %v@%v\n", in.GetUsername(), in.GetHostname())

	// Crea un canale con cui comprendere se l'username è valido o meno
	validityCh := make(chan bool)

	// Invia i metadati relativi al nuovo utente alla goroutine dedicata all'aggiornamento
	// del datastore
	s.newMemberCh <- &exchangeInformationFromUpdaterAndHandler{
		clientInfo:        in,
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
		return status.Errorf(codes.InvalidArgument, "Username already in use!")
	}
}

func NewRegistrationServer(size uint16) *RegistrationServer {
	return &RegistrationServer{
		newMemberCh:       make(chan *exchangeInformationFromUpdaterAndHandler),
		memberInformation: []*exchangeInformationFromUpdaterAndHandler{},
		numberOfClients:   size,
	}
}

func ServeSignRequests(exposedPort uint16, regServer *RegistrationServer, grpcServer *grpc.Server, wg *sync.WaitGroup) {
	defer wg.Done()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", exposedPort))
	if err != nil {
		panic(err)
	}

	proto.RegisterRegistrationServer(grpcServer, regServer)
	grpcServer.Serve(lis)
}
