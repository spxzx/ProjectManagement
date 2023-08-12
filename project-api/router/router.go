package router

import (
	"github.com/gin-gonic/gin"
)

// Router 路由接口，用于分模块实现具体的路由处理
type Router interface {
	Route(r *gin.Engine)
}

var routers []Router

func Register(ro ...Router) {
	routers = append(routers, ro...)
}

// RegisterRouter 路由注册器
type RegisterRouter struct {
}

// New 路由注册器实例生成
func New() *RegisterRouter {
	return &RegisterRouter{}
}

// Route 调用实现的接口
//
//   - ro 为传入的具体路由结构
func (*RegisterRouter) Route(ro Router, r *gin.Engine) {
	ro.Route(r)
}

// InitRouter 路由初始化
func InitRouter(r *gin.Engine) {
	for _, ro := range routers {
		ro.Route(r)
	}
}
