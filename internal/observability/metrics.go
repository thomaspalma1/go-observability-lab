package observability

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// httpRequestsTotal conta o total de requisições HTTP recebidas,
// segmentado por método, rota e status de resposta.
var httpRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total de requisições HTTP recebidas",
	},
	[]string{"method", "path", "status"},
)

// httpRequestDuration mede a duração de cada requisição HTTP,
// segmentado por método e rota.
var httpRequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duração das requisições HTTP em segundos",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path"},
)

// Metrics é o middleware que registra as métricas de cada requisição.
// Deve ser registrado depois do RequestID e antes dos handlers de rota.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // processa a requisição

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = "unmatched"
		}

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}
