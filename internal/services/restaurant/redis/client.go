package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *redis.Client
	Ctx context.Context
}

func New(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Client{
		RDB: rdb,
		Ctx: context.Background(),
	}
}

func (c *Client) SetOrderStatus(orderID string, status string) error {
	return c.RDB.Set(c.Ctx, "order:"+orderID+":status", status, 0).Err()
}
