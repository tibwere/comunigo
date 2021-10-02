package grpchandler

import (
	"fmt"
	"log"
	"net"
	"sync"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

type RegistrationServer struct {
	proto.UnimplementedRegistrationServer
	updateSync      []*chan bool
	clientChan      chan *proto.ClientInfo
	chatGroup       []*proto.ClientInfo
	numberOfClients uint16
}

func (s *RegistrationServer) UpdateMembers(grpcServer *grpc.Server, seqAddr string, seqPort uint16, wg *sync.WaitGroup) {
	defer wg.Done()

	for uint16(len(s.chatGroup)) < s.numberOfClients {
		member := <-s.clientChan
		s.chatGroup = append(s.chatGroup, member)
	}

	InitializeSequencer(seqAddr, seqPort, s.chatGroup)

	// Invia un messaggio di sincronizzazione a tutte
	// le goroutine per inviare la risposta al client
	for _, ch := range s.updateSync {
		*ch <- true
	}

	// stop del server poiché è stato completato il gruppo di multicast
	grpcServer.GracefulStop()
}

func (s *RegistrationServer) Sign(in *proto.ClientInfo, stream proto.Registration_SignServer) error {
	log.Printf("Received request from: %v@%v\n", in.GetUsername(), in.GetHostname())

	// Invia il nuovo utente alla goroutine dedicata all'aggiornamento
	// del datastore
	s.clientChan <- in

	// Crea un canale con cui sincronizzarsi sull'invio delle risposte
	syncChan := make(chan bool)
	s.updateSync = append(s.updateSync, &syncChan)

	// Aspetta il messaggio di sincronizzazione per restituire le info al chiamante
	<-syncChan

	for _, member := range s.chatGroup {
		if err := stream.Send(member); err != nil {
			return err
		}
	}
	return nil
}

func NewRegistrationServer(size uint16) *RegistrationServer {
	return &RegistrationServer{
		updateSync:      []*chan bool{},
		clientChan:      make(chan *proto.ClientInfo),
		chatGroup:       []*proto.ClientInfo{},
		numberOfClients: size,
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
