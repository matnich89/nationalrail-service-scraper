package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	NationalRailApiKey string
	RedisAddress       string
	TrainIdQueueName   string
}

func Load() (*Config, error) {
	// Load .env file only if not in k8s environment
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	return &Config{
		NationalRailApiKey: getEnv("NATIONAL_RAIL_API_KEY"),
		RedisAddress:       getEnv("REDIS_ADDRESS"),
		TrainIdQueueName:   getEnv("TRAIN_ID_QUEUE_NAME"),
	}, nil
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}
