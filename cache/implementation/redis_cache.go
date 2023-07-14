package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping the Redis server to ensure the connection is successful
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Get(key string, v interface{}) error {
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key '%s' not found in cache", key)
		}
		return fmt.Errorf("failed to get key '%s' from cache: %v", key, err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal cache data for key '%s': %v", key, err)
	}

	return nil
}

func (c *RedisCache) Set(key string, v interface{}, expiration time.Duration) error {
	ctx := context.Background()

	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal data for key '%s': %v", key, err)
	}

	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key '%s' in cache: %v", key, err)
	}

	return nil
}

func (c *RedisCache) Delete(key string) error {
	ctx := context.Background()
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key '%s' from cache: %v", key, err)
	}

	return nil
}
