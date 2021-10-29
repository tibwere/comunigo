package grpchandler

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/grpc"
)

// "Metodo della classe SequencerServer" eseguito da una goroutine che permette di inizializzare le connessioni
// ai peer e al termine di spengere il server gRPC per la ricezione dei peer dal nodo di registrazione
func (s *SequencerServer) StartupConnectionWithPeers(ctx context.Context, fromRegToSeqGRPCserver *grpc.Server) error {
	errCh := make(chan error)

	for i := 0; i < int(s.chatGroupSize); i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case currentMember := <-s.memberCh:
			s.connections[currentMember.GetAddress()] = make(chan *proto.SequencerMessage)

			go func() {
				if err := s.sendBackMessages(ctx, currentMember.GetAddress()); err != nil {
					errCh <- err
				}
			}()
		}
	}

	// tutte le connessioni sono state aperte quindi è possibile
	// stoppare il server GRPC
	fromRegToSeqGRPCserver.GracefulStop()

	// costruzione del messaggio di errore da restituire
	errMsg := ""
	for addr := range s.connections {
		errMsg += fmt.Sprintf("Handler for: %v->%v, ", addr, <-errCh)
	}
	// rimuove l'ulitmo ", "
	errMsg = errMsg[:len(errMsg)-2]

	return fmt.Errorf(errMsg)

}

// "Metodo della classe SequencerServer" che permette l'invio
// dei messaggi 'preparati' dalla goroutine ad-hoc all'i-esimo peer connesso
func (s *SequencerServer) sendBackMessages(ctx context.Context, addr string) error {
	conn, err := grpc.Dial(
		fmt.Sprintf("%v:%v", addr, s.port),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := proto.NewComunigoClient(conn)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case newMessage := <-s.connections[addr]:
			_, err = c.SendFromSequencerToPeer(context.Background(), newMessage)
			if err != nil {
				return err
			}
		}
	}
}

// "Metodo della classe SequencerServer" per l'implementazione della RPC SendFromPeerToSequencer server-side
func (s *SequencerServer) SendFromPeerToSequencer(ctx context.Context, in *proto.RawMessage) (*empty.Empty, error) {
	log.Printf("Received '%v' from %v\n", in.GetBody(), in.GetFrom())
	s.seqCh <- in
	return &empty.Empty{}, nil
}

// "Metodo della classe SequencerServer" eseguito da una goroutine che riceve
// i messaggi tramite un canale in cui le procedure remote li inseriscono
// li ordina aggiungendo un timestamp e li inoltra alle procedure dedicate
// all'invio
func (s *SequencerServer) OrderMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("signal caught")
		case unordered := <-s.seqCh:
			ordered := &proto.SequencerMessage{
				Timestamp: s.sequenceNumber,
				From:      unordered.GetFrom(),
				Body:      unordered.GetBody(),
			}
			s.sequenceNumber++

			for _, ch := range s.connections {
				ch <- ordered
			}
		}
	}
}

// "Metodo della classe SequencerServer" che inizializza il server gRPC per la ricezione
// dei messaggi dai peer partecipanti al gruppo di multicast
func (s *SequencerServer) ServePeers(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	proto.RegisterComunigoServer(grpcServer, s)
	go grpcServer.Serve(lis)

	<-ctx.Done()
	log.Println("Message receiver from sequencer shutdown")
	grpcServer.GracefulStop()
	return fmt.Errorf("signal caught")
}
