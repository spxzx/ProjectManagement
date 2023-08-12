package main

import (
	"github.com/gin-gonic/gin"
	srv "github.com/spxzx/project-common"
	"github.com/spxzx/project-project/config"
	"github.com/spxzx/project-project/router"
	"github.com/spxzx/project-project/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
)

func main() {
	r := gin.Default()

	tp, tpErr := tracing.JaegerTraceProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	r.Use(otelgin.Middleware("project-project"))

	grpc := router.RegisterGRPC()
	router.RegisterEtcdServer()
	router.InitUserRpcClient()

	//kwc := config.InitKafkaWriter()

	srv.Run(
		r,
		config.Conf.Server.Port,
		config.Conf.Server.Name,
		func() {
			grpc.Stop()
			//kwc()
		},
	)
}
