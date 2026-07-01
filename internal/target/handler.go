package target

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra as rotas do serviço alvo, que recebe as requisições
// disparadas pelo load runner.
func RegisterRoutes(router *gin.Engine) {
	router.GET("/target/ping", handlePing)
}

// handlePing responde a requisições de teste de carga
//
//	@Summary	Endpoint alvo para testes de carga
//	@Tags		target
//	@Produce	json
//	@Success	200	{object}	map[string]bool
//	@Router		/target/ping [get]
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"pong": true})
}
