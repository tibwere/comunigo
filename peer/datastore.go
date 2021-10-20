package peer

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
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

func GetMessages(ds *redis.Client, key string) ([]string, error) {
	var messages []string
	ctx := context.Background()

	messages, err := ds.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return messages, err
	}

	log.Printf("Found %v messages into redis to deliver to frontend (key: %v)\n", len(messages), key)

	return messages, nil
}
