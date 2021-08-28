package main

import (
	"context"
	"log"
	"net"

	idl "gitlab.com/tibwere/comunigo/greetings/encoding/proto"
	"google.golang.org/grpc"
)

const (
	port = ":29314"
)

type server struct {
	idl.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *idl.HelloRequest) (*idl.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &idl.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	idl.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
