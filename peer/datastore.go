package peer

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gitlab.com/tibwere/comunigo/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

func InitDatastore(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:6379", addr),
		Password: "",
		DB:       0,
	})
}

func InsertMessage(ds *redis.Client, rootKey string, receivedMessage *proto.OrderedMessage) error {
	mOpt := &protojson.MarshalOptions{
		Multiline:       false,
		EmitUnpopulated: true,
	}

	byteMessage, err := mOpt.Marshal(receivedMessage)
	if err != nil {
		return err
	} else {

		return ds.Set(
			context.Background(),
			fmt.Sprintf("%v-%v", rootKey, receivedMessage.GetID()),
			string(byteMessage),
			0,
		).Err()
	}
}

func GetMessages(ds *redis.Client, rootKey string, startID uint64) ([]string, error) {
	currentIndex := startID
	var messages []string
	ctx := context.Background()

	for {
		mes, err := ds.Get(
			ctx,
			fmt.Sprintf("%v-%v", rootKey, currentIndex),
		).Result()

		if err == redis.Nil {
			return messages, nil
		} else if err != nil {
			return nil, err
		} else {
			messages = append(messages, mes)
			currentIndex++
		}
	}
}
