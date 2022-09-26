package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sloth/config"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis(c config.Conf) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:         c.Redis.Host + ":" + c.Redis.Port,
		Password:     c.Redis.Password,
		DB:           c.Redis.DB,
		ReadTimeout:  c.Redis.ReadTimeout,
		WriteTimeout: c.Redis.WriteTimeout,
		DialTimeout:  c.Redis.DialTimeout,
	})
	_, err := Rdb.Ping(Ctx).Result()

	return err
}
