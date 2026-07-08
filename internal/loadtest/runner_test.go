package loadtest

import (
	"context"
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

	result := Run(context.Background(), cfg)

	if result.TotalRequests.Load() == 0 {
		t.Fatal("expected at least one request to be sent, but none were executed")
	}

	if result.Failed.Load() > 0 {
		t.Errorf("expected zero failures against a healthy target, but got %d", result.Failed.Load())
	}

	if result.Successful.Load() != result.TotalRequests.Load() {
		t.Errorf(
			"expected all requests to succeed: total=%d successful=%d",
			result.TotalRequests.Load(),
			result.Successful.Load(),
		)
	}
}
