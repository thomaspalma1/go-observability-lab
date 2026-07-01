package main

import (
	"log"

	"github.com/gin-gonic/gin"
//	"github.com/thomaspalma1/go-observability-lab/internal/loadtest"
	"github.com/thomaspalma1/go-observability-lab/internal/target"
)

func main() {
	router := gin.Default()

	// Rota de teste de vida, só pra validar que o servidor está de pé.
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Registra as rotas do endpoint "alvo" (quem recebe a carga)
	target.RegisterRoutes(router)

	// Registra as rotas do disparador de carga
	//loadtest.RegisterRoutes(router)


	if err := router.Run(":8082"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
