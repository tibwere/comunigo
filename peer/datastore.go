package peer

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// "Metodo della classe Status" che permette di inizializzare
// la connessione al datastore (redis)
func (s *Status) initDatastore(addr string) {
	s.datastore = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:6379", addr),
		Password: "",
		DB:       0,
	})
}

// "Metodo della classe Status" che permette di effettuare l'RPUSH del messaggio
// all'interno del datastore
func (s *Status) RPUSHMessage(message protoreflect.ProtoMessage) error {
	enc := &protojson.MarshalOptions{
		Multiline:       false,
		EmitUnpopulated: true,
	}

	byteMessage, err := enc.Marshal(message)
	if err != nil {
		return err
	} else {
		val := string(byteMessage)
		log.Printf("RPush into redis at key %v val: %v\n", s.currentUsername, val)
		return s.datastore.RPush(context.Background(), s.currentUsername, val).Err()
	}
}

// "Metodo della classe Status" che semplicemente serve da prologo per i
// metodi GetMessage*()
func (s *Status) getMessagesPrologue() ([]string, error) {
	ctx := context.Background()
	return s.datastore.LRange(ctx, s.currentUsername, 0, -1).Result()
}

// "Metodo della classe Status" che permette di effettuare il retrieve
// dei messaggi ricevuti dal sequencer
func (s *Status) GetMessagesSEQ() ([]*proto.SequencerMessage, error) {
	messages := []*proto.SequencerMessage{}

	rawMessages, err := s.getMessagesPrologue()
	if err != nil {
		return messages, err
	}

	log.Printf(
		"Found %v messages into redis to deliver to frontend (key: %v)\n",
		len(rawMessages),
		s.currentUsername,
	)

	for _, raw := range rawMessages {
		mess := &proto.SequencerMessage{}
		protojson.Unmarshal([]byte(raw), mess)
		messages = append(messages, mess)
	}

	return messages, nil
}

// "Metodo della classe Status" che permette di effettuare il retrieve
// dei messaggi dagli altri peer in modo totalmente ordinato
func (s *Status) GetMessagesSC() ([]*proto.ScalarClockMessage, error) {
	messages := []*proto.ScalarClockMessage{}

	rawMessages, err := s.getMessagesPrologue()
	if err != nil {
		return messages, err
	}

	log.Printf(
		"Found %v messages into redis to deliver to frontend (key: %v)\n",
		len(rawMessages),
		s.currentUsername,
	)

	for _, raw := range rawMessages {
		mess := &proto.ScalarClockMessage{}
		protojson.Unmarshal([]byte(raw), mess)
		messages = append(messages, mess)
	}

	return messages, nil
}

// "Metodo della classe Status" che permette di effettuare il retrieve
// dei messaggi dagli altri peer in modo causalmente ordinato
func (s *Status) GetMessagesVC() ([]*proto.VectorialClockMessage, error) {
	messages := []*proto.VectorialClockMessage{}

	rawMessages, err := s.getMessagesPrologue()
	if err != nil {
		return messages, err
	}

	log.Printf(
		"Found %v messages into redis to deliver to frontend (key: %v)\n",
		len(rawMessages),
		s.currentUsername,
	)

	for _, raw := range rawMessages {
		mess := &proto.VectorialClockMessage{}
		protojson.Unmarshal([]byte(raw), mess)
		messages = append(messages, mess)
	}

	return messages, nil
}
