package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedis(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func (r *RedisClient) SetOrderStatus(orderID string, status string) error {
	return r.Client.Set(r.Ctx, "order:"+orderID+":status", status, 0).Err()
}

func (r *RedisClient) GetOrderStatus(orderID string) (string, error) {
	return r.Client.Get(r.Ctx, "order:"+orderID+":status").Result()
}
