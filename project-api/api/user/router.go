package user

import (
	"github.com/gin-gonic/gin"
	"github.com/spxzx/project-api/api/mid"
	"github.com/spxzx/project-api/api/rpc"
	"github.com/spxzx/project-api/router"
)

type RouterUser struct{}

func init() {
	router.Register(&RouterUser{})
}

func (*RouterUser) Route(r *gin.Engine) {
	// 初始化grpc客户端连接
	rpc.InitUserRpcClient()
	h := New()

	r.POST("/project/login/getCaptcha", h.getCaptcha)
	r.POST("/project/login/register", h.register)
	r.POST("/project/login", h.login)

	o := r.Group("/project/organization", mid.TokenVerify())
	o.POST("/_getOrgList", h.getMyOrgList)
}
