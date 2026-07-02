package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/thomaspalma1/go-observability-lab/cmd/api/docs"
	"github.com/thomaspalma1/go-observability-lab/internal/health"
	"github.com/thomaspalma1/go-observability-lab/internal/loadtest"
	"github.com/thomaspalma1/go-observability-lab/internal/observability"
	"github.com/thomaspalma1/go-observability-lab/internal/target"
)

// @title			Go Observability Lab API
// @version		1.0
// @description	Load runner e alvo simulado para estudo de observabilidade
// @host			localhost:8082
// @BasePath		/
func main() {
	logger := observability.NewLogger()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(observability.RequestID())
	router.Use(observability.RequestLogger(logger))
	router.Use(observability.Metrics())

	health.RegisterRoutes(router, loadtest.ActiveTests)
	target.RegisterRoutes(router)
	loadtest.RegisterRoutes(router)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8082"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
