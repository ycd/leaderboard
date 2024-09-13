package storage

import (
	"context"
	"log"
	"testing"
)

func TestRedis(t *testing.T) {
	client, err := NewRedisConnection()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	client.Set(context.Background(), "test", "test", 0)

	val, err := client.Get(context.Background(), "test").Result()
	if err != nil {
		log.Fatalf("Could not get value from Redis: %v", err)
	}

	log.Printf("Value from Redis: %v", val)
}
