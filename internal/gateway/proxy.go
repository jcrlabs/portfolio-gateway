package gateway

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jonathanCaamano/portfolio-gateway/internal/config"
	"github.com/jonathanCaamano/portfolio-gateway/internal/metrics"
)

// NewMux builds the ServeMux with all routes and middleware chain.
func NewMux(routes []config.Route, cfg *config.Config, log *slog.Logger, m *metrics.Metrics) http.Handler {
	mux := http.NewServeMux()

	for _, route := range routes {
		r := route
		proxy := newReverseProxy(r.Backend, log)
		handler := chain(proxy,
			withRequestID,
			withLogger(log),
			withCORS(cfg.AllowOrigins),
			withRateLimit(cfg.RateLimitRPS),
			withRecover(log),
			withMetrics(m, r.Name),
		)
		// Strip the prefix so backends receive the path without /api/{service}
		mux.Handle(r.Prefix+"/", http.StripPrefix(r.Prefix, handler))
	}

	mux.Handle("/healthz", healthHandler(routes, log))
	mux.Handle("/metrics", m.Handler())

	return mux
}

func newReverseProxy(backend string, log *slog.Logger) http.Handler {
	target, err := url.Parse(backend)
	if err != nil {
		panic(fmt.Sprintf("invalid backend URL %q: %v", backend, err))
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error("proxy error", "backend", backend, "path", r.URL.Path, "err", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"bad gateway"}`, http.StatusBadGateway)
	}
	return proxy
}
