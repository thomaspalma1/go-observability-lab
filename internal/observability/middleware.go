package observability

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"
const requestIDKey = "request_id"

// RequestID gera um ID único por requisição, expõe no header de resposta
// e disponibiliza no contexto do Gin para outros middlewares/handlers usarem.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewString()
		c.Set(requestIDKey, id)
		c.Writer.Header().Set(requestIDHeader, id)
		c.Next()
	}
}

// RequestLogger loga cada requisição em JSON, incluindo o request_id,
// método, rota, status e duração.
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // processa a requisição

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
