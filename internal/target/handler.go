package target

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra as rotas do serviço alvo, que recebe as requisições
// disparadas pelo load runner.
func RegisterRoutes(router *gin.Engine) {
	router.GET("/target/ping", handlePing)
}

// handlePing responde a requisições de teste de carga. Aceita um parâmetro
// opcional delay_ms, usado para simular lentidão e testar SLOs/alertas.
//
//	@Summary	Endpoint alvo para testes de carga
//	@Tags		target
//	@Produce	json
//	@Param		delay_ms	query		int	false	"Atraso artificial em milissegundos"
//	@Success	200			{object}	map[string]bool
//	@Router		/target/ping [get]
func handlePing(c *gin.Context) {
	if delayParam := c.Query("delay_ms"); delayParam != "" {
		if delayMs, err := strconv.Atoi(delayParam); err == nil && delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
	}

	c.JSON(http.StatusOK, gin.H{"pong": true})
}
