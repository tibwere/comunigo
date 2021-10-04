package peer

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"gitlab.com/tibwere/comunigo/proto"
)

func InitDatastore(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:6379", addr),
		Password: "",
		DB:       0,
	})
}

func InsertMessage(ds *redis.Client, key string, value *proto.OrderedMessage) error {
	jsonMessage, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ds.LPush(key, jsonMessage).Err()
}
