package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// maxConcurrentTests define o limite de testes de carga simultâneos que
// consideramos "saudável". Acima disso, o serviço se reporta como não-pronto.
const maxConcurrentTests = 3

// activeTestsFunc é o tipo da função usada para consultar testes ativos.
// Definido como tipo para facilitar testes e evitar acoplamento direto ao
// pacote loadtest.
type activeTestsFunc func() int64

// RegisterRoutes registra as rotas de liveness e readiness.
//
// getActiveTests é injetado de fora (main.go) para evitar que o pacote
// health precise importar o pacote loadtest diretamente.
func RegisterRoutes(router *gin.Engine, getActiveTests activeTestsFunc) {
	router.GET("/healthz", handleLiveness)
	router.GET("/readyz", handleReadiness(getActiveTests))
}

// handleLiveness confirma que o processo está de pé, sem checar dependências.
//
//	@Summary	Liveness probe
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/healthz [get]
func handleLiveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// handleReadiness verifica se o serviço está apto a receber tráfego,
// considerando o número de testes de carga rodando simultaneamente.
//
//	@Summary	Readiness probe
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Failure	503	{object}	map[string]interface{}
//	@Router		/readyz [get]
func handleReadiness(getActiveTests activeTestsFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		active := getActiveTests()

		if active > maxConcurrentTests {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":       "not_ready",
				"reason":       "too many concurrent load tests",
				"active_tests": active,
				"max_allowed":  maxConcurrentTests,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":       "ready",
			"active_tests": active,
		})
	}
}
