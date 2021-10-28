package reg

import (
	"context"
	"fmt"
	"io"
	"log"

	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// "Metodo della classe ToRegisterGRPCHandler" che permette di effettuare il retrieve
// dei peer connessi a partire dallo stream gRPC della risposta del nodo di registrazione
//
// Restituisce inoltre true nel caso in cui l'username inserito è valido
func (h *ToRegisterGRPCHandler) getOtherMembers(currUser string, stream proto.Registration_SignClient) (bool, error) {
	for {
		member, err := stream.Recv()

		if err == io.EOF {
			return true, nil
		}

		if err != nil {
			errStatus, _ := status.FromError(err)
			if codes.InvalidArgument == errStatus.Code() {
				h.peerStatus.PushIntoFrontendBackendChannel(errStatus.Message())
				return false, nil
			} else {
				return true, err
			}
		}

		if currUser != member.GetUsername() {
			h.peerStatus.InsertNewMember(member)
		}
	}
}

// "Metodo della classe ToRegisterGRPCHandler" che permette di invocare il servizio RPC
// esposto dal nodo di registrazione per richiedere la registrazione utilizzando l'username
// correntemente inserito dall'utente e ricevuto tramite canale
func (h *ToRegisterGRPCHandler) SignToRegister(ctx context.Context) error {
	var currUser string

	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", h.registerAddr, h.registerPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("Connection established to: %v:%v\n", h.registerAddr, h.registerPort)

	c := proto.NewRegistrationClient(conn)

	for {
		select {
		case <-ctx.Done():
			log.Println("Registration client shutdown")
			return fmt.Errorf("signal caught")

		case currUser = <-h.peerStatus.GetFromFrontendBackendChannel():
			stream, err := c.Sign(context.Background(), &proto.NewUser{
				Username: currUser,
			})
			if err != nil {
				return err
			}

			validUsername, err := h.getOtherMembers(currUser, stream)
			if err != nil {
				return err
			}

			if validUsername {
				h.peerStatus.PushIntoFrontendBackendChannel("SUCCESS")
				return h.peerStatus.SetUsername(currUser)
			}
		}
	}
}
