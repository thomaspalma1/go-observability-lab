package loadtest

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Result guarda os números agregados de uma execução de teste de carga.
type Result struct {
	TotalRequests int64
	Successful    int64
	Failed        int64

	mu        sync.Mutex
	latencies []time.Duration
}

// Config define os parâmetros de uma execução de teste de carga.
type Config struct {
	TargetURL         string
	RequestsPerSecond int
	Duration          time.Duration
}

// Run dispara requisições HTTP contra TargetURL, tentando manter o ritmo de
// RequestsPerSecond, durante Duration. Bloqueia até o teste terminar.
func Run(cfg Config) *Result {
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
			return result
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				fireRequest(client, cfg.TargetURL, result)
			}()
		}
	}
}

func fireRequest(client *http.Client, url string, result *Result) {
	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)

	result.mu.Lock()
	result.latencies = append(result.latencies, elapsed)
	result.mu.Unlock()

	atomic.AddInt64(&result.TotalRequests, 1)

	if err != nil || resp.StatusCode >= http.StatusInternalServerError {
		atomic.AddInt64(&result.Failed, 1)
		return
	}
	defer resp.Body.Close()
	atomic.AddInt64(&result.Successful, 1)
}
