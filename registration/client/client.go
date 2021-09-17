package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	pb "gitlab.com/tibwere/comunigo/registration/proto"
	"google.golang.org/grpc"
)

var usernamePtr = flag.String("u", "", "Username for comuniGO chat group")
var serverPtr = flag.String("s", "registration-server", "Registration server")

func main() {
	flag.Parse()
	fmt.Printf("Tentativo di connessione a %s:2929\n", *serverPtr)

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:2929", *serverPtr),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewRegistrationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream, err := c.Sign(ctx, &pb.Request{
		Username: *usernamePtr,
	})

	fmt.Println("In attesa degli altri partecipanti ...")
	if err != nil {
		panic(err)
	}

	counter := 1
	for {
		info, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Membro n.%v: %v\n", counter, info)
		counter++
	}
}
