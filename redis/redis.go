package redis

import (
	"context"
	"encoding/json"
	"fmt"
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

func NewRedisClient(queueName string, redisAddress string) (IRedisClient, error) {
	r := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	ctx := context.Background()
	_, err := r.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisClient{queueName: queueName, client: r}, nil
}

func (r *RedisClient) PushToQueue(ctx context.Context, id model.DepartingTrainId) {

	departingIdJSON, err := json.Marshal(id)
	if err != nil {
		log.Printf("error serializing train %s to JSON: %v", departingIdJSON, err)
	}

	err = r.client.RPush(ctx, r.queueName, departingIdJSON).Err()
	if err != nil {
		log.Printf("error adding train %s to Redis queue: %v", id.ID, err)
	}
}
