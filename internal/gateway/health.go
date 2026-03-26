package gateway

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jonathanCaamano/portfolio-gateway/internal/config"
)

type backendStatus struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
}

func healthHandler(routes []config.Route, log *slog.Logger) http.Handler {
	client := &http.Client{Timeout: 3 * time.Second}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		results := make([]backendStatus, 0, len(routes))
		overall := "ok"

		for _, route := range routes {
			start := time.Now()
			status := "ok"
			resp, err := client.Get(route.Backend + "/api/health")
			ms := time.Since(start).Milliseconds()
			if err != nil || (resp != nil && resp.StatusCode >= 400) {
				status = "degraded"
				overall = "degraded"
			}
			if resp != nil {
				resp.Body.Close()
			}
			results = append(results, backendStatus{
				Name:      route.Name,
				Status:    status,
				LatencyMs: ms,
			})
		}

		code := http.StatusOK
		if overall != "ok" {
			code = http.StatusServiceUnavailable
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":   overall,
			"backends": results,
		})
		log.Info("healthz", "status", overall)
	})
}
