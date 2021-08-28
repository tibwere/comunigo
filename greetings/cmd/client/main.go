package main

import (
	"context"
	"log"
	"os"
	"time"

	"gitlab.com/tibwere/comunigo/greetings/encoding/proto"
	"google.golang.org/grpc"
)

const (
	socket = "localhost:29314"
)

func main() {
	if len(os.Args) <= 1 {
		// TODO: Aggiungere stampa usage
		log.Fatalf("Missing command line parameters")
	}

	conn, err := grpc.Dial(socket, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: os.Args[1]})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
}
