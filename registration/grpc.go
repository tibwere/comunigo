package registration

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/go-redis/redis"
	pb "gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

type RegistrationServer struct {
	pb.UnimplementedRegistrationServer
	updateSync      []*chan bool
	clientChan      chan *pb.ClientInfo
	chatGroup       *pb.ChatGroupInfo
	numberOfClients uint16
}

func (s *RegistrationServer) saveMembersOnRedis() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "registration_ds:6379",
		Password: "",
		DB:       0,
	})
	jsonInfos, err := json.Marshal(s.chatGroup.Infos)
	if err != nil {
		panic(err)
	}
	redisClient.Set("chat-group", jsonInfos, 0)
}

func (s *RegistrationServer) UpdateMembers() {

	for uint16(len(s.chatGroup.Infos)) < s.numberOfClients {
		fmt.Printf("Sono qui %v - %v\n", uint16(len(s.chatGroup.Infos)), s.numberOfClients)
		member := <-s.clientChan
		s.chatGroup.Infos = append(s.chatGroup.Infos, member)
	}

	fmt.Println("Sono uscito")
	// Invia un messaggio di sincronizzazione a tutte
	// le goroutine per inviare la risposta al client
	for _, ch := range s.updateSync {
		fmt.Println("Sto inviando i sync")
		*ch <- true
	}
}

func (s *RegistrationServer) Sign(ctx context.Context, in *pb.ClientInfo) (*pb.ChatGroupInfo, error) {
	fmt.Printf("Received request from: %v@%v\n", in.GetUsername(), in.GetHostname())

	// Invia il nuovo utente alla goroutine dedicata all'aggiornamento
	// del datastore
	s.clientChan <- in

	// Crea un canale con cui sincronizzarsi sull'invio delle risposte
	syncChan := make(chan bool)
	s.updateSync = append(s.updateSync, &syncChan)

	// Aspetta il messaggio di sincronizzazione per restituire le info al chiamante
	<-syncChan
	return s.chatGroup, nil
}

func NewRegistrationServer(size uint16) *RegistrationServer {
	return &RegistrationServer{
		updateSync:      []*chan bool{},
		clientChan:      make(chan *pb.ClientInfo),
		chatGroup:       &pb.ChatGroupInfo{},
		numberOfClients: size,
	}
}

func ServeSignRequests(exposedPort uint16, regServer *RegistrationServer) {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", exposedPort))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRegistrationServer(grpcServer, regServer)

	grpcServer.Serve(lis)
}
