package storage

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

// CacheGet retrieves the cached content from redis.
func (s *Storage) CacheGet(ctx context.Context, key string) (interface{}, error) {
	resp, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var raw interface{}
	if err := json.Unmarshal([]byte(resp), &raw); err != nil {
		log.Printf("got error while unmarshalling: %v", err)
	}

	return raw, nil
}

// CacheSet caches the given content with the 60 seconds of TTL
func (s *Storage) CacheSet(ctx context.Context, key string, leaderboard []byte) error {
	_, err := s.cache.Set(ctx, key, leaderboard, 60*time.Second).Result()
	if err != nil {
		return err
	}

	return nil
}
