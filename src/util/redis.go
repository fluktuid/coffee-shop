package util

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()

type Client struct {
	rdb *redis.Client
}

func NewClient(addr string) *Client {
	logrus.Info(addr)
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Client{rdb}
}

func (c *Client) LPush(key, val string) error {
	err := c.rdb.LPush(ctx, key, val).Err()
	return err
}

func (c *Client) BRPop(key string) (string, error) {
	val, err := c.rdb.BRPop(ctx, 10+time.Second, key).Result()
	if len(val) >= 2 {
		return val[1], err
	}
	return "", err
}

func (c *Client) Set(key, val string, expiration time.Duration) error {
	err := c.rdb.Set(ctx, key, val, expiration).Err()
	return err
}

func (c *Client) GetDel(key string) (string, error) {
	val, err := c.rdb.GetDel(ctx, key).Result()
	return val, err
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
