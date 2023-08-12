package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spxzx/project-user/config"
	"time"
)

type RedisCache struct {
	rdb *redis.Client
}

var Rc *RedisCache

func init() {
	rdb := redis.NewClient(config.Conf.InitRedisOptions())
	Rc = &RedisCache{rdb}
}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	return rc.rdb.Set(ctx, key, value, expire).Err()
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return rc.rdb.Get(ctx, key).Result()
}
