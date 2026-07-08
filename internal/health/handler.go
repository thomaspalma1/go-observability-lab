package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// maxConcurrentTests defines the maximum number of concurrent load tests
// considered healthy. If this threshold is exceeded, the service reports
// itself as not ready.
const maxConcurrentTests = 3

// activeTestsFunc represents the function used to retrieve the current number
// of active load tests. It is defined as a type to simplify testing and avoid
// coupling this package directly to the loadtest package.
type activeTestsFunc func() int64

// RegisterRoutes registers the liveness and readiness endpoints.
//
// The getActiveTests function is injected from main.go so that the health
// package does not need to import the loadtest package directly.
func RegisterRoutes(router *gin.Engine, getActiveTests activeTestsFunc) {
	router.GET("/healthz", handleLiveness)
	router.GET("/readyz", handleReadiness(getActiveTests))
}

// handleLiveness confirms that the process is running without checking any
// external dependencies.
//
//	@Summary	Liveness probe
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/healthz [get]
func handleLiveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// handleReadiness verifies whether the service is ready to receive traffic
// based on the number of concurrently running load tests.
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
