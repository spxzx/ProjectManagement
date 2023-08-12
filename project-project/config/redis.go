package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spxzx/project-project/internal/dao"
)

func (c *config) ReConnRedis() {
	rdb := redis.NewClient(c.InitRedisOptions())
	rc := &dao.RedisCache{Rdb: rdb}
	dao.Rc = rc
}
