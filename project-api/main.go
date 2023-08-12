package main

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spxzx/project-api/config"
	"github.com/spxzx/project-api/tracing"
	srv "github.com/spxzx/project-common"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"net/http"

	// _ 不引用，只是做个索引让init函数进行初始化
	_ "github.com/spxzx/project-api/api"
	"github.com/spxzx/project-api/router"
)

func main() {
	r := gin.Default()

	//r.Use(cors.Cors())

	tp, tpErr := tracing.JaegerTraceProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	r.Use(otelgin.Middleware("project-api"))

	r.StaticFS("/upload", http.Dir("upload"))
	router.InitRouter(r)

	// 开启 pprof 默认访问路径 /debug/pprof
	pprof.Register(r)
	// goroutine 内存泄露测试

	srv.Run(
		r,
		config.Conf.Server.Port,
		config.Conf.Server.Name,
		nil, // 这里是连接gRPC服务所以不需要
	)
}
