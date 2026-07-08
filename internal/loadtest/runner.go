package loadtest

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("go-observability-lab/loadtest")

// Result stores the aggregated metrics collected during a load test execution.
type Result struct {
	TotalRequests atomic.Int64
	Successful    atomic.Int64
	Failed        atomic.Int64

	mu        sync.Mutex
	latencies []time.Duration
}

// Config defines the parameters of a load test execution.
type Config struct {
	TargetURL         string
	RequestsPerSecond int
	Duration          time.Duration
}

// Run sends HTTP requests to the target URL while attempting to maintain the
// configured requests-per-second rate for the specified duration.
//
// The function blocks until the load test has completed.
func Run(ctx context.Context, cfg Config) *Result {
	ctx, span := tracer.Start(ctx, "loadtest.run",
		trace.WithAttributes(
			attribute.String("target.url", cfg.TargetURL),
			attribute.Int("requests_per_second", cfg.RequestsPerSecond),
			attribute.String("duration", cfg.Duration.String()),
		),
	)
	defer span.End()

	result := &Result{}

	interval := time.Second / time.Duration(cfg.RequestsPerSecond)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	stop := time.After(cfg.Duration)
	var wg sync.WaitGroup

	client := &http.Client{Timeout: 5 * time.Second}

	for {
		select {
		case <-stop:
			wg.Wait()
			span.SetAttributes(
				attribute.Int64("requests.total", result.TotalRequests.Load()),
				attribute.Int64("requests.successful", result.Successful.Load()),
				attribute.Int64("requests.failed", result.Failed.Load()),
			)
			return result
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				fireRequest(ctx, client, cfg.TargetURL, result)
			}()
		}
	}
}

// fireRequest performs a single HTTP request and records its outcome,
// latency, and tracing information.
func fireRequest(ctx context.Context, client *http.Client, url string, result *Result) {
	_, span := tracer.Start(ctx, "loadtest.fire_request",
		trace.WithAttributes(
			attribute.String("http.url", url),
		),
	)
	defer span.End()

	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)

	span.SetAttributes(attribute.String("http.duration", elapsed.String()))

	result.mu.Lock()
	result.latencies = append(result.latencies, elapsed)
	result.mu.Unlock()

	result.TotalRequests.Add(1)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		result.Failed.Add(1)
		return
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	if resp.StatusCode >= http.StatusInternalServerError {
		span.SetStatus(codes.Error, "server error")
		result.Failed.Add(1)
		return
	}

	result.Successful.Add(1)
}
