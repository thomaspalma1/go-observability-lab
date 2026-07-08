package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	_ "github.com/thomaspalma1/go-observability-lab/cmd/api/docs"
	"github.com/thomaspalma1/go-observability-lab/internal/health"
	"github.com/thomaspalma1/go-observability-lab/internal/loadtest"
	"github.com/thomaspalma1/go-observability-lab/internal/observability"
	"github.com/thomaspalma1/go-observability-lab/internal/target"
)

// @title			Go Observability Lab API
// @version		1.0
// @description	Load generator and simulated target for observability studies.
// @host			localhost:8082
// @BasePath		/
func main() {
	ctx := context.Background()

	shutdown, err := observability.InitTracer(ctx, "go-observability-lab", "jaeger:4317")
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	defer shutdown(ctx)

	logger := observability.NewLogger()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(observability.RequestID())
	router.Use(observability.RequestLogger(logger))
	router.Use(observability.Metrics())
	router.Use(otelgin.Middleware("go-observability-lab"))

	health.RegisterRoutes(router, loadtest.ActiveTests)
	target.RegisterRoutes(router)
	loadtest.RegisterRoutes(router)
	observability.RegisterPprof(router)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8082"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
