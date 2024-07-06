package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"trainstats-scraper/model"
)

type IRedisClient interface {
	PushToQueue(ctx context.Context, id model.DepartingTrainId)
}

type RedisClient struct {
	queueName string
	client    *redis.Client
}

func NewRedisClient(queueName string, redisAddress string) IRedisClient {
	r := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	return &RedisClient{queueName: queueName, client: r}
}

func (r *RedisClient) PushToQueue(ctx context.Context, id model.DepartingTrainId) {

	departingIdJSON, err := json.Marshal(id)
	if err != nil {
		log.Printf("error serializing train %s to JSON: %v", departingIdJSON, err)
	}

	err = r.client.RPush(ctx, "train_id_queue", departingIdJSON).Err()
	if err != nil {
		log.Printf("error adding train %s to Redis queue: %v", id.ID, err)
	}
}
