// Package per la gestione della logica applicativa
// basata sullo scambio dei messaggi gRPC tra nodo di
// nodo di registrazione e peers/sequencer
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

// In ottica OO, oggetto che rappresenta l'informazione
// scambiata fra le goroutine spawnate automaticamente dal
// server gRPC per l'esecuzione della procedura remota
// e della goroutine responsaibile della logica di registrazione
type exchangeInformationFromUpdaterAndHandler struct {
	clientInfo        *proto.PeerInfo
	isUsernameValidCh chan bool
}

// In ottica OO, oggetto che rappresenta il server di registrazione
type RegistrationServer struct {
	proto.UnimplementedRegistrationServer
	newMemberCh       chan *exchangeInformationFromUpdaterAndHandler
	memberInformation []*exchangeInformationFromUpdaterAndHandler
	numberOfClients   uint16
}

// "Costruttore" dell'oggetto RegistrationServer
func NewRegistrationServer(size uint16) *RegistrationServer {
	return &RegistrationServer{
		newMemberCh:       make(chan *exchangeInformationFromUpdaterAndHandler),
		memberInformation: []*exchangeInformationFromUpdaterAndHandler{},
		numberOfClients:   size,
	}
}

// "Metodo della classe RegistrationServer" che a partire da un'array
// di elementi del tipo exchangeInformationFromUpdaterAndHandler restituisce
// la lista di membri del gruppo di multicast che si sono appena registrati
func (s *RegistrationServer) getChatGroupMembers() []*proto.PeerInfo {
	var members []*proto.PeerInfo

	for _, info := range s.memberInformation {
		members = append(members, info.clientInfo)
	}

	return members
}

// "Metodo della classe RegistrationServer" che a partire da un'array
// di elementi del tipo exchangeInformationFromUpdaterAndHandler restituisce
// la lista dei canali su cui le varie procedure remote stanno attendendo
// per capire se poter consegnare la lista al chiamante oppure restituire un errore
func (s *RegistrationServer) getValidityList() []chan bool {
	var validityList []chan bool

	for _, info := range s.memberInformation {
		validityList = append(validityList, info.isUsernameValidCh)
	}

	return validityList
}

// "Metodo della classe RegistrationServer" che permette di verificare se un dato username
// inserito da un peer è valido o meno
//
// la regola di validità si basa sull'unicità dell'username all'interno della lista
func (s *RegistrationServer) isValidUsername(username string) bool {
	for _, member := range s.memberInformation {
		if member.clientInfo.GetUsername() == username {
			return false
		}
	}
	return true
}

// "Metodo della classe RegistrationServer" eseguito da una goroutine che:
//
// 1) Riceve dalle varie procedure remote le informazioni
//
// 2) Verifica la validità degli username richiesti
//
// 3) Inizializza il sequencer
//
// 4) Sblocca le procedure remote per far si che consegnino il risultato al chiamante
//
// 5) Stoppa il server di registrazione
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

// "Metodo della classe RegistrationServer" per l'implementazione della RPC Sign server-side
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

// "Metodo della classe P2PHandler" che inizializza il server gRPC per la ricezione
// dei messaggi dai peer che vogliono registrarsi al gruppo di multicast
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
