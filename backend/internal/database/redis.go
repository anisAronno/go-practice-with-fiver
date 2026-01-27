package database

import (
	"context"
	"fmt"
	"gofiver/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the Redis connection
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient creates a new Redis connection
func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &RedisClient{Client: client}, nil
}

// Set stores a value with expiration
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Delete removes a key
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
