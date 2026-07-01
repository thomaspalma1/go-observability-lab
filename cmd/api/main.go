package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/thomaspalma1/go-observability-lab/internal/loadtest"
	"github.com/thomaspalma1/go-observability-lab/internal/target"
)

func main() {
	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	target.RegisterRoutes(router)
	loadtest.RegisterRoutes(router)

	if err := router.Run(":8082"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
