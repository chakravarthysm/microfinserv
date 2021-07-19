package common

import (
	"errors"

	"github.com/go-redis/redis"
)

type Client struct {
	redis *redis.Client
}

func NewRedisClient() (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Client{redis: client}, nil
}

func (c *Client) IsBlacklisted(jwt string) error {
	m, err := c.redis.SMembersMap("blacklist").Result()
	if err != nil {
		return err
	}

	if _, ok := m[jwt]; ok {
		return errors.New("token has been blacklisted")
	}

	return nil
}

func (c *Client) AddToBlacklist(jwt string) error {
	_, err := c.redis.SAdd("blacklist", jwt).Result()
	if err != nil {
		return err
	}
	return err
}
