package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	pb "gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

func getRemoteServerFromEnv() (string, error) {
	serverAddr, isPresent := os.LookupEnv("COMUNIGO_RHOST")
	if !isPresent {
		return "", errors.New("please specify server address")
	}

	return serverAddr, nil
}

var usernamePtr = flag.String("u", "", "Username for comuniGO chat group")

func main() {

	flag.Parse()

	serverAddr, err := getRemoteServerFromEnv()
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:2929", serverAddr),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("In attesa degli altri partecipanti ...")

	c := pb.NewRegistrationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	res, err := c.Sign(ctx, &pb.ClientInfo{
		Username: *usernamePtr,
		Hostname: os.Getenv("HOSTNAME"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Modalità: %v\n", res.Tos)
	for i, client := range res.GetInfos() {
		fmt.Printf("Membro %v: %v@%v\n", i, client.GetUsername(), client.GetHostname())
	}

}
