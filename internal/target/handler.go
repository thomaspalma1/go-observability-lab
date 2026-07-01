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

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"pong": true})
}
