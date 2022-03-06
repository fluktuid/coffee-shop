package util

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Client struct {
	rdb *redis.Client
}

func NewClient(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Client{rdb}
}

func (c *Client) LPush(key, val string) error {
	err := c.rdb.LPush(ctx, key, val, 0).Err()
	return err
}

func (c *Client) BRPop(key string) (string, error) {
	val, err := c.rdb.BRPop(ctx, 10+time.Second, key).Result()
	if len(val) > 0 {
		return val[0], err
	}
	return "", err
}

func (c *Client) Set(key, val string, expiration time.Duration) error {
	err := c.rdb.Set(ctx, key, val, expiration).Err()
	return err
}

func (c *Client) GetDel(key string) (string, error) {
	val, err := c.rdb.Do(ctx, "MULTI", "GET", key, "DEL", key, "EXEC").Result()
	return val.(string), err
}
func (c *Client) AnyPattern(pattern string) (string, error) {
	resp, err := c.rdb.Keys(ctx, pattern).Result()
	for _, e := range resp {
		res, err := c.rdb.GetDel(ctx, e).Result()
		if res != "" {
			return res, err
		}
	}
	return "", err
}
