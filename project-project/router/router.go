package router

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/spxzx/project-common/discovery"
	"github.com/spxzx/project-common/logs"
	"github.com/spxzx/project-grpc/account/account"
	"github.com/spxzx/project-grpc/auth/auth"
	"github.com/spxzx/project-grpc/department/department"
	"github.com/spxzx/project-grpc/menu/menu"
	"github.com/spxzx/project-grpc/project/project"
	"github.com/spxzx/project-grpc/task/task"
	"github.com/spxzx/project-project/config"
	"github.com/spxzx/project-project/internal/rpc"
	accountService "github.com/spxzx/project-project/pkg/service/account"
	authService "github.com/spxzx/project-project/pkg/service/auth"
	departmentService "github.com/spxzx/project-project/pkg/service/department"
	menuService "github.com/spxzx/project-project/pkg/service/menu"
	projectService "github.com/spxzx/project-project/pkg/service/project"
	taskService "github.com/spxzx/project-project/pkg/service/task"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
)

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterGRPC() *grpc.Server {
	c := gRPCConfig{
		Addr: config.Conf.GRPC.Addr,
		RegisterFunc: func(s *grpc.Server) {
			project.RegisterProjectServiceServer(s, projectService.New())
			task.RegisterTaskServiceServer(s, taskService.New())
			account.RegisterAccountServiceServer(s, accountService.New())
			department.RegisterDepartmentServiceServer(s, departmentService.New())
			auth.RegisterAuthServiceServer(s, authService.New())
			menu.RegisterMenuServiceServer(s, menuService.New())
		},
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		otelgrpc.UnaryServerInterceptor(),
		//interceptor.New().CacheInterceptor(),
	)))
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		zap.L().Error("cannot listen, cause by: " + err.Error())
	}
	go func() {
		log.Printf("grpc server started on: %s \n", c.Addr)
		if err := s.Serve(lis); err != nil {
			zap.L().Error("server start error, cause by: " + err.Error())
			return
		}
	}()
	return s
}

func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.Conf.Etcd.Addrs, logs.Log)
	resolver.Register(etcdRegister)
	r := discovery.NewRegister(config.Conf.Etcd.Addrs, logs.Log)
	if _, err := r.Register(discovery.Server{
		Name:    config.Conf.GRPC.Name,
		Addr:    config.Conf.GRPC.Addr,
		Version: config.Conf.GRPC.Version,
		Weight:  config.Conf.GRPC.Weight,
	}, 2); err != nil {
		log.Fatalln(err)
	}
}

func InitUserRpcClient() {
	rpc.InitUserRpcClient()
}
