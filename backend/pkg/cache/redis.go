package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client.
// Supports both direct host:port and Upstash-style URL connections.
func NewRedisClient(ctx context.Context, redisURL, host, port, password string, db int) (*redis.Client, error) {
	var client *redis.Client

	if redisURL != "" {
		// Parse URL (e.g., rediss://default:xxx@region.upstash.io:6379)
		opts, err := redis.ParseURL(redisURL)
		if err != nil {
			return nil, fmt.Errorf("parse redis URL: %w", err)
		}
		client = redis.NewClient(opts)
	} else {
		// Direct connection
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       db,
		})
	}

	// Verify connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
