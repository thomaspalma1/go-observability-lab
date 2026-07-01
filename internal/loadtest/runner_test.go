package loadtest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRun_FiresRequestsAgainstTarget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := Config{
		TargetURL:         server.URL,
		RequestsPerSecond: 10,
		Duration:          500 * time.Millisecond,
	}

	result := Run(cfg)

	if result.TotalRequests == 0 {
		t.Fatal("esperava pelo menos uma requisição disparada, mas não disparou nenhuma")
	}

	if result.Failed > 0 {
		t.Errorf("esperava zero falhas contra um alvo saudável, mas teve %d", result.Failed)
	}

	if result.Successful != result.TotalRequests {
		t.Errorf("esperava que todas as requisições fossem bem-sucedidas: total=%d sucesso=%d",
			result.TotalRequests, result.Successful)
	}
}
