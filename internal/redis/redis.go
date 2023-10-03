package redis

import (
	"context"
	"net/url"
	"os"

	"github.com/redis/go-redis/v9"
)

type Connection struct {
	client *redis.Client
}

func Connect() (Connection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0, // use default DB
	})

	return Connection{client: client}, nil
}

func (c *Connection) Get(key string) (url.URL, error) {
	ctx := context.Background()
	res, err := c.client.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}

	url, err := url.Parse(res)
	if err != nil {
		panic(err)
	}

	return *url, nil
}

func (c *Connection) Set(key string, value url.URL) error {
	ctx := context.Background()
	err := c.client.Set(ctx, key, value.String(), 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}
