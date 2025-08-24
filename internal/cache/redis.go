package cache

import (
	"context"

	"github.com/qullDev/book_API/internal/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, Password: cfg.RedisPassword, DB: cfg.RedisDB})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err

	}

	return rdb, nil

}
