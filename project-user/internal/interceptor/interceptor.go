package interceptor

import (
	"context"
	"encoding/json"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-grpc/user/login"
	"github.com/spxzx/project-user/internal/dao"
	"github.com/spxzx/project-user/internal/repo"
	"google.golang.org/grpc"
	"time"
)

type CacheInterceptor struct {
	cache repo.Cache
	cMap  map[string]any
}

func New() *CacheInterceptor {
	cMap := make(map[string]any)
	cMap["/login.service.LoginService/FindMemInfoById"] = &login.OrgResponse{}
	cMap["/login.service.LoginService/GetOrgList"] = &login.Member{}
	return &CacheInterceptor{
		cache: dao.Rc,
		cMap:  cMap,
	}
}

func (ci *CacheInterceptor) CacheInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		respType := ci.cMap[info.FullMethod]
		if respType == nil {
			return handler(ctx, req)
		}
		// 先查询是否有缓存 有=>直接返回 没有=>请求再存入缓存
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		reqJson, _ := json.Marshal(req)
		cKey := encrypts.Md5(string(reqJson))

		respJson, err := ci.cache.Get(c, info.FullMethod+"::"+cKey)
		if respJson != "" {
			_ = json.Unmarshal([]byte(respJson), &respType)
			return respType, nil
		}

		resp, err = handler(ctx, req)
		bytes, _ := json.Marshal(resp)

		_ = ci.cache.Put(c, info.FullMethod+"::"+cKey, string(bytes), time.Minute)
		return
	}
}
