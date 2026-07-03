package loadtest

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// runRequest é o corpo esperado em POST /load-test/run
type runRequest struct {
	TargetURL         string `json:"target_url" binding:"required"`
	RequestsPerSecond int    `json:"requests_per_second" binding:"required,gt=0"`
	DurationSeconds   int    `json:"duration_seconds" binding:"required,gt=0"`
}

var (
	resultsMu sync.RWMutex
	results   = make(map[string]*Result)

	activeTests atomic.Int64
)

// ActiveTests retorna quantos testes de carga estão rodando neste momento.
// Usado pelo readiness check para refletir a saúde interna do processo.
func ActiveTests() int64 {
	return activeTests.Load()
}

// RegisterRoutes registra as rotas do load runner.
func RegisterRoutes(router *gin.Engine) {
	router.POST("/load-test/run", handleRun)
	router.GET("/load-test/:id/results", handleResults)
}

// handleRun inicia um novo teste de carga
//
//	@Summary	Inicia um teste de carga
//	@Tags		load-test
//	@Accept		json
//	@Produce	json
//	@Param		request	body		runRequest	true	"Configuração do teste"
//	@Success	202		{object}	map[string]string
//	@Failure	400		{object}	map[string]string
//	@Router		/load-test/run [post]
func handleRun(c *gin.Context) {
	var req runRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	testID := uuid.NewString()

	cfg := Config{
		TargetURL:         req.TargetURL,
		RequestsPerSecond: req.RequestsPerSecond,
		Duration:          time.Duration(req.DurationSeconds) * time.Second,
	}

	// WithoutCancel preserva o trace_id da requisição original, mas remove
	// o cancelamento automático - necessário porque o teste de carga
	// continua rodando em background, depois da resposta HTTP já ter voltado.
	ctx := context.WithoutCancel(c.Request.Context())

	activeTests.Add(1)

	go func() {
		defer activeTests.Add(-1)

		result := Run(ctx, cfg)
		resultsMu.Lock()
		results[testID] = result
		resultsMu.Unlock()
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"test_id": testID,
		"status":  "started",
	})
}

// handleResults consulta o resultado de um teste de carga
//
//	@Summary	Consulta resultado de um teste
//	@Tags		load-test
//	@Produce	json
//	@Param		id	path		string	true	"Test ID"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]string
//	@Router		/load-test/{id}/results [get]
func handleResults(c *gin.Context) {
	testID := c.Param("id")

	resultsMu.RLock()
	result, found := results[testID]
	resultsMu.RUnlock()

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "test not found or still running"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"test_id":        testID,
		"total_requests": result.TotalRequests.Load(),
		"successful":     result.Successful.Load(),
		"failed":         result.Failed.Load(),
	})
}
