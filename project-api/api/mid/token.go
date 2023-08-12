package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/spxzx/project-api/api/rpc"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/r"
	"github.com/spxzx/project-grpc/user/login"
)

func TokenVerify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1.从header中获取token
		token := ctx.GetHeader("Authorization")
		// 2.调用user模块进行token认证
		tokenResp, err := rpc.LoginServiceClient.TokenVerify(ctx, &login.TokenVerifyRequest{Token: token, Ip: GetIp(ctx)})
		// 3.处理结果 放入gin上下文/返回未登录
		if err != nil {
			e := errs.ParseGrpcError(err)
			r.Fail(ctx, int(e.Code), e.Msg)
			ctx.Abort()
			return
		}
		ctx.Set("memberId", tokenResp.Member.Id)
		ctx.Set("memberName", tokenResp.Member.Name)
		ctx.Set("organizationCode", tokenResp.Member.OrganizationCode)
		ctx.Next()
	}
}

func GetIp(c *gin.Context) (ip string) {
	ip = c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return
}
