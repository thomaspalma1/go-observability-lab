package observability

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// httpRequestsTotal counts the total number of HTTP requests received,
// grouped by method, route, and response status.
var httpRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests received",
	},
	[]string{"method", "path", "status"},
)

// httpRequestDuration measures the duration of each HTTP request,
// grouped by method and route.
var httpRequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path"},
)

// Metrics is a middleware that records metrics for every HTTP request.
//
// It should be registered after the RequestID middleware and before the
// route handlers.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // Process the request.

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = "unmatched"
		}

		httpRequestsTotal.
			WithLabelValues(c.Request.Method, path, status).
			Inc()

		httpRequestDuration.
			WithLabelValues(c.Request.Method, path).
			Observe(duration)
	}
}
