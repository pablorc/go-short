package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Connection struct {
	client *redis.Client
}

func Connect() (Connection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return Connection{client: client}, nil
}

func (c *Connection) Get(key string) (string, error) {
	ctx := context.Background()
	res, err := c.client.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return res, nil
}

func (c *Connection) Set(key, value string) error {
	ctx := context.Background()
	err := c.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
	return nil
}
