package redis

import (
	"context"
	"fmt"

	"github.com/jevvonn/readora-backend/config"
	"github.com/redis/go-redis/v9"
)

func New() *redis.Client {
	conf := config.Load()

	addr := fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return rdb
}
