package redis

import "github.com/go-redis/redis/v8"

type Client struct {
	*redis.Client
}
