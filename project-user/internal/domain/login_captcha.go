package domain

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-user/internal/dao"
	"github.com/spxzx/project-user/internal/repo"
	"github.com/spxzx/project-user/pkg/data"
	"go.uber.org/zap"
	"time"
)

type LoginCaptchaDomain struct {
	cache repo.Cache
}

func (c *LoginCaptchaDomain) PutCaptcha(key, value string, exp time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.cache.Put(ctx, key, value, exp)
}

func (c *LoginCaptchaDomain) CheckCaptcha(mobile, captcha string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	getCaptcha, err := c.cache.Get(ctx, data.RegisterKey+mobile)
	if err == redis.Nil {
		return errs.GrpcError(data.CaptchaNotExist)
	}
	if err != nil {
		zap.L().Error("register redis get error, cause by: ", zap.Error(err))
		return errs.GrpcError(data.RedisError)
	}
	if getCaptcha != captcha {
		return errs.GrpcError(data.CaptchaError)
	}
	return nil
}

func NewLoginCaptchaDomain() *LoginCaptchaDomain {
	return &LoginCaptchaDomain{
		cache: dao.Rc,
	}
}
