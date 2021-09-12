package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"gitlab.com/tibwere/comunigo/registry/entity"
	pb "gitlab.com/tibwere/comunigo/registry/proto"
	"google.golang.org/grpc"
)

type configuration struct {
	Hosts []string `json:"hosts"`
	Tos   string   `json:"tos"`
	Port  uint32   `json:"port"`
}

type registryServer struct {
	pb.UnimplementedRegistryServer
}

var cl = entity.SafeClientList{}
var cfg configuration
var canStart chan bool

func (s *registryServer) Sign(ctx context.Context, in *pb.ClientInfo) (*pb.TypeOfService, error) {

	for _, host := range cfg.Hosts {
		if host == in.GetHostname() {
			cl.Add(in.GetUsername(), in.GetHostname())

			fmt.Println(cl.HowMany(), len(cfg.Hosts))

			if cl.HowMany() == len(cfg.Hosts) {
				canStart <- true
			}

			log.Printf(
				"Inserted new user: %s@%s (Local port: %d)",
				in.GetUsername(), in.GetHostname(),
				in.GetLocalPort(),
			)
			return &pb.TypeOfService{Type: cfg.Tos}, nil
		}
	}

	log.Printf("Host not allowed tried to connect (%v)\n", in.GetHostname())
	return nil, errors.New("host not allowed")

}

func loadConfig(config *configuration) error {
	configFile, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer configFile.Close()

	data, err := io.ReadAll(configFile)
	if err != nil {
		return err
	}

	json.Unmarshal(data, config)

	return nil
}

func startChat() {
	active := <-canStart

	if active {
		fmt.Println("Ora si può iniziare")
	}
}

func main() {
	canStart = make(chan bool)
	err := loadConfig(&cfg)
	if err != nil {
		log.Fatalf("Load config failed (%x)", err)
	}

	fmt.Println("Hosts: ")
	for i, host := range cfg.Hosts {
		fmt.Printf("\t%d) %s\n", i+1, host)
	}
	fmt.Printf("Type of service: %s\nExposed port: %d\n", cfg.Tos, cfg.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}

	go startChat()

	grpcServer := grpc.NewServer()
	pb.RegisterRegistryServer(grpcServer, &registryServer{})
	grpcServer.Serve(lis)
}
