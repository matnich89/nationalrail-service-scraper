package config

import (
	"errors"
	"fmt"
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
		NationalRailApiKey: os.Getenv("NATIONAL_RAIL_API_KEY"),
	}, nil
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", errors.New(fmt.Sprintf("Environment variable %s is not set", key))
}
