package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisCache struct {
	Rdb *redis.Client
}

var Rc *RedisCache

//func init() {
//	rdb := redis.NewClient(config.Conf.InitRedisOptions())
//	Rc = &RedisCache{rdb}
//}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	return rc.Rdb.Set(ctx, key, value, expire).Err()
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return rc.Rdb.Get(ctx, key).Result()
}
