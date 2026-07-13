package redis

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)


type Client struct {
    rdb *redis.Client
}

func New(url string) (*Client, error) {
    opts, err := redis.ParseURL(url)
    if err != nil {
        return nil, err
    }
    rdb := redis.NewClient(opts)

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    if err := rdb.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    return &Client{rdb: rdb}, nil
}

func (c *Client) Close() error {
    return c.rdb.Close()
}


func (c *Client) Increment(ctx context.Context, key string, ttl time.Duration) (int64, error) {
    pipe := c.rdb.Pipeline()
    incr := pipe.Incr(ctx, key)
    pipe.Expire(ctx, key, ttl)
    _, err := pipe.Exec(ctx)
    if err != nil {
        return 0, err
    }
    return incr.Val(), nil
}

func (c *Client) Ping(ctx context.Context) error {
    return c.rdb.Ping(ctx).Err()
}