package cache

import (
	"context"
	"time"

	"github.com/notOliveira/onde-tem/internal/core/ports"
	"github.com/redis/go-redis/v9"
)

var _ ports.Cache = (*ValkeyClient)(nil)

type ValkeyClient struct {
	client *redis.Client
}

func NewValkeyClient(addr string) *ValkeyClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &ValkeyClient{client: rdb}
}

func (v *ValkeyClient) Get(ctx context.Context, key string) (string, error) {
	return v.client.Get(ctx, key).Result()
}

func (v *ValkeyClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return v.client.Set(ctx, key, value, ttl).Err()
}

func (v *ValkeyClient) Delete(ctx context.Context, key string) error {
	return v.client.Del(ctx, key).Err()
}
