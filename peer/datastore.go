package peer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gitlab.com/tibwere/comunigo/proto"
)

func InitDatastore(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:6379", addr),
		Password: "",
		DB:       0,
	})
}

func InsertMessage(ds *redis.Client, rootKey string, receivedMessage *proto.OrderedMessage) error {
	jsonMessage, err := json.Marshal(receivedMessage)
	if err != nil {
		return err
	} else {
		fmt.Printf("Sto inserendo la chiave %v-%v", rootKey, receivedMessage.GetID())

		return ds.Set(
			context.Background(),
			fmt.Sprintf("%v-%v", rootKey, receivedMessage.GetID()),
			jsonMessage,
			0,
		).Err()
	}
}

func GetMessages(ds *redis.Client, rootKey string, startID uint64) ([]*proto.OrderedMessage, error) {
	currentIndex := startID
	var messages []*proto.OrderedMessage
	var mes *proto.OrderedMessage
	ctx := context.Background()

	for {
		jsonMes, err := ds.Get(
			ctx,
			fmt.Sprintf("%v-%v", rootKey, currentIndex),
		).Result()

		if err == redis.Nil {
			return messages, nil
		} else if err != nil {
			return nil, err
		} else {
			json.Unmarshal([]byte(jsonMes), &mes)
			fmt.Println(mes)
			messages = append(messages, mes)
			currentIndex++
		}
	}
}
