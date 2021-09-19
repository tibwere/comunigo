package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	pb "gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

type RegistrationServer struct {
	pb.UnimplementedRegistrationServer
	updateSync      []*chan bool
	clientChan      chan *pb.ClientInfo
	chatGroup       *pb.ChatGroupInfo
	numberOfClients int
}

func (s *RegistrationServer) UpdateMembers() {

	for len(s.chatGroup.Infos) < s.numberOfClients {
		member := <-s.clientChan
		s.chatGroup.Infos = append(s.chatGroup.Infos, member)
	}

	// Memorizza l'insieme dei client connessi su redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	jsonInfos, err := json.Marshal(s.chatGroup.Infos)
	if err != nil {
		panic(err)
	}
	redisClient.Set("chat-group", jsonInfos, 0)

	// Invia un messaggio di sincronizzazione a tutte
	// le goroutine per inviare la risposta al client
	for _, ch := range s.updateSync {
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

func setChatGroupSize() (int, error) {
	chosenSize, isPresent := os.LookupEnv("SIZE")
	if isPresent {
		return strconv.Atoi(chosenSize)
	} else {
		return 3, nil
	}
}

func main() {

	lis, err := net.Listen("tcp", ":2929")
	if err != nil {
		panic(err)
	}

	size, err := setChatGroupSize()
	if err != nil {
		panic(err)
	}

	fmt.Printf("In ascolto sulla porta 2929 (Dimensione del gruppo: %v)\n", size)

	grpcServer := grpc.NewServer()
	regServer := &RegistrationServer{
		updateSync:      []*chan bool{},
		clientChan:      make(chan *pb.ClientInfo),
		chatGroup:       &pb.ChatGroupInfo{},
		numberOfClients: size,
	}

	regServer.chatGroup.Tos = pb.TypeOfService_SEQUENCER

	go regServer.UpdateMembers()
	pb.RegisterRegistrationServer(grpcServer, regServer)
	grpcServer.Serve(lis)
}
