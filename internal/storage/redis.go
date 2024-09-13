package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection() (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	log.Printf("Connecting to Redis at %s with password %s", addr, redisPassword)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword,
		Username: "",
		DB:       redisDB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %v", err)
	}

	log.Printf("Connected to Redis")

	return client, nil
}
