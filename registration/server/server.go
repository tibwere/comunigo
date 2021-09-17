package main

import (
	"errors"
	"flag"
	"fmt"
	"net"

	pb "gitlab.com/tibwere/comunigo/registration/proto"
	"google.golang.org/grpc"
)

type Configuration struct {
	Hosts []string `json:"hosts"`
}

type MulticastMember struct {
	Info            *pb.Request
	ResponseChannel chan string
}

type RegistrationServer struct {
	pb.UnimplementedRegistrationServer
	syncro          chan *MulticastMember
	actualClients   map[string]chan string
	numberOfClients int
}

func (s *RegistrationServer) UpdateMembers() {

	for len(s.actualClients) < s.numberOfClients {
		member := <-s.syncro
		s.actualClients[member.Info.GetUsername()] = member.ResponseChannel
	}

	s.SendMemberList()
}

func (s *RegistrationServer) SendMemberList() {
	for clientToBeUpdated, ch := range s.actualClients {
		for member := range s.actualClients {
			if clientToBeUpdated != member {
				ch <- member
			}
		}

		// Send empty username to signal end of comunication
		ch <- ""
	}
}

func (s *RegistrationServer) Sign(in *pb.Request, stream pb.Registration_SignServer) error {
	fmt.Printf("Received request from: %v\n", in.GetUsername())

	if in.GetUsername() == "" {
		return errors.New("username must be not empty")
	}

	responses := make(chan string)
	s.syncro <- &MulticastMember{
		Info:            in,
		ResponseChannel: responses,
	}

	for {
		info := <-responses
		if info == "" {
			return nil
		}
		if err := stream.Send(&pb.Response{Username: info}); err != nil {
			return err
		}
	}
}

var sizePtr = flag.Int("n", 3, "Number of clients allowed")

func main() {

	flag.Parse()

	lis, err := net.Listen("tcp", ":2929")
	if err != nil {
		panic(err)
	}

	fmt.Printf("In ascolto sulla porta 2929 (Dimensione del gruppo: %v)\n", *sizePtr)

	grpcServer := grpc.NewServer()
	regServer := &RegistrationServer{
		syncro:          make(chan *MulticastMember),
		actualClients:   make(map[string]chan string),
		numberOfClients: *sizePtr,
	}

	go regServer.UpdateMembers()
	pb.RegisterRegistrationServer(grpcServer, regServer)
	grpcServer.Serve(lis)
}
