package rpc

import (
	"github.com/spxzx/project-api/config"
	"github.com/spxzx/project-common/discovery"
	"github.com/spxzx/project-common/logs"
	"github.com/spxzx/project-grpc/user/login"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"os"
)

var LoginServiceClient login.LoginServiceClient

func InitUserRpcClient() {
	etcdRegister := discovery.NewResolver(config.Conf.Etcd.Addrs, logs.Log)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial(
		etcdRegister.Scheme()+":///project-user", // :/// 后面的名字要和 grpc配置 中的 name 一致
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		zap.L().Error("did not connect, cause by: " + err.Error())
		os.Exit(1)
	}
	LoginServiceClient = login.NewLoginServiceClient(conn)
}
