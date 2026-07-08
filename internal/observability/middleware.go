package observability

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"
const requestIDKey = "request_id"

// RequestID generates a unique ID for each incoming request, exposes it in the
// response header, and stores it in the Gin context so it can be used by other
// middlewares and handlers.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewString()
		c.Set(requestIDKey, id)
		c.Writer.Header().Set(requestIDHeader, id)
		c.Next()
	}
}

// RequestLogger logs every HTTP request as a JSON record, including the
// request ID, method, route, response status, request duration, and client IP.
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // Process the request.

		duration := time.Since(start)
		requestID, _ := c.Get(requestIDKey)

		logger.Info("http_request",
			slog.String("request_id", requestID.(string)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Int("status", c.Writer.Status()),
			slog.String("duration", duration.String()),
			slog.String("client_ip", c.ClientIP()),
		)
	}
}
