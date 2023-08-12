package rpc

import (
	"github.com/spxzx/project-common/discovery"
	"github.com/spxzx/project-common/logs"
	"github.com/spxzx/project-grpc/account/account"
	"github.com/spxzx/project-grpc/auth/auth"
	"github.com/spxzx/project-grpc/department/department"
	"github.com/spxzx/project-grpc/menu/menu"
	"github.com/spxzx/project-grpc/project/project"
	"github.com/spxzx/project-grpc/task/task"
	"github.com/spxzx/project-project/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"os"
)

var ProjectServiceClient project.ProjectServiceClient
var TaskServiceClient task.TaskServiceClient
var AccountServiceClient account.AccountServiceClient
var DepartmentServiceClient department.DepartmentServiceClient
var AuthServiceClient auth.AuthServiceClient
var MenuServiceClient menu.MenuServiceClient

func InitProjectRpcClient() {
	etcdRegister := discovery.NewResolver(config.Conf.Etcd.Addrs, logs.Log)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial(
		etcdRegister.Scheme()+":///project-project", // :/// 后面的名字要和 grpc配置 中的 name 一致
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		zap.L().Error("did not connect, cause by: " + err.Error())
		os.Exit(1)
	}
	ProjectServiceClient = project.NewProjectServiceClient(conn)
	TaskServiceClient = task.NewTaskServiceClient(conn)
	AccountServiceClient = account.NewAccountServiceClient(conn)
	DepartmentServiceClient = department.NewDepartmentServiceClient(conn)
	AuthServiceClient = auth.NewAuthServiceClient(conn)
	MenuServiceClient = menu.NewMenuServiceClient(conn)
}
