package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "gitlab.com/tibwere/comunigo/registry/proto"
	"google.golang.org/grpc"
)

var address = flag.String("a", "", "Registry server address")
var rport = flag.Uint("s", 8080, "Remote port")
var lport = flag.Uint("w", 8081, "Local port")
var username = flag.String("u", "", "Username")
var hostname = flag.String("m", "", "Hostname")

func signToRegistry(socket string, username string, hostname string, lport uint32) (string, error) {
	conn, err := grpc.Dial(socket, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return "", nil
	}
	defer conn.Close()

	c := pb.NewRegistryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	res, err := c.Sign(ctx, &pb.ClientInfo{
		Username:  username,
		Hostname:  hostname,
		LocalPort: lport,
	})

	return res.GetType(), err
}

// func waitForOthers(socket string) ([]entity.Client, error) {
// 	_, err := net.Listen("tcp", socket)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//grpcServer := grpc.NewServer()
// 	//pb.RegisterStarterServer(grpcServer, &registryServer{})
// 	// grpcServer.Serve(lis)

// 	return nil, nil
// }

func main() {
	flag.Parse()

	if *username == "" || *hostname == "" || *address == "" {
		flag.Usage()
		return
	}

	mod, err := signToRegistry(
		fmt.Sprintf("%s:%d", *address, *rport),
		*username,
		*hostname,
		uint32(*lport),
	)
	if err != nil {
		log.Fatalf("Unable to sign to registry (%v)", err)
	}
	log.Printf("Successfully signed.\nModality: %s\nWaiting for other partecipants ...", mod)

}
