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

func InitDatastore(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:6379", addr),
		Password: "",
		DB:       0,
	})
}

func RPUSHMessage(ds *redis.Client, key string, message protoreflect.ProtoMessage) error {
	enc := &protojson.MarshalOptions{
		Multiline:       false,
		EmitUnpopulated: true,
	}

	byteMessage, err := enc.Marshal(message)
	if err != nil {
		return err
	} else {
		val := string(byteMessage)
		log.Printf("RPush into redis at key %v val: %v\n", key, val)
		return ds.RPush(context.Background(), key, val).Err()
	}
}

func getMessagesPrologue(ds *redis.Client, key string) ([]string, error) {
	ctx := context.Background()
	return ds.LRange(ctx, key, 0, -1).Result()
}

func GetMessagesSEQ(ds *redis.Client, key string) ([]*proto.SequencerMessage, error) {
	messages := []*proto.SequencerMessage{}

	rawMessages, err := getMessagesPrologue(ds, key)
	if err != nil {
		return messages, err
	}

	log.Printf("Found %v messages into redis to deliver to frontend (key: %v)\n", len(rawMessages), key)

	for _, raw := range rawMessages {
		mess := &proto.SequencerMessage{}
		protojson.Unmarshal([]byte(raw), mess)
		messages = append(messages, mess)
	}

	return messages, nil
}

func GetMessagesSC(ds *redis.Client, key string) ([]*proto.ScalarClockMessage, error) {
	messages := []*proto.ScalarClockMessage{}

	rawMessages, err := getMessagesPrologue(ds, key)
	if err != nil {
		return messages, err
	}

	log.Printf("Found %v messages into redis to deliver to frontend (key: %v)\n", len(rawMessages), key)

	for _, raw := range rawMessages {
		mess := &proto.ScalarClockMessage{}
		protojson.Unmarshal([]byte(raw), mess)
		messages = append(messages, mess)
	}

	return messages, nil
}

func GetMessagesVC(ds *redis.Client, key string) ([]*proto.VectorialClockMessage, error) {
	messages := []*proto.VectorialClockMessage{}

	rawMessages, err := getMessagesPrologue(ds, key)
	if err != nil {
		return messages, err
	}

	log.Printf("Found %v messages into redis to deliver to frontend (key: %v)\n", len(rawMessages), key)

	for _, raw := range rawMessages {
		mess := &proto.VectorialClockMessage{}
		protojson.Unmarshal([]byte(raw), mess)
		messages = append(messages, mess)
	}

	return messages, nil
}
