package project

import (
	"github.com/gin-gonic/gin"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/r"
	"go.uber.org/zap"
	"strings"
)

var ignores = []string{
	"project/login/register",
	"project/login",
	"project/login/getCaptcha",
	"project/organization",
	"project/auth/apply",
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zap.L().Info("开始做授权认证")
		//当用户登录认证通过，获取到用户信息，查询用户权限所拥有的节点信息
		//根据请求的uri路径 进行匹配
		uri := ctx.Request.RequestURI
		a := NewAuth()
		nodes, err := a.GetAuthNodes(ctx)
		if err != nil {
			e := errs.ParseGrpcError(err)
			r.Fail(ctx, int(e.Code), e.Msg)
			ctx.Abort()
			return
		}
		for _, v := range ignores {
			if strings.Contains(uri, v) {
				ctx.Next()
				return
			}
		}
		for _, v := range nodes {
			if strings.Contains(uri, v) {
				ctx.Next()
				return
			}
		}
		r.Fail(ctx, 403, "无操作权限")
		ctx.Abort()
	}
}
