package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	latency  *prometheus.HistogramVec
	requests *prometheus.CounterVec
}

func New() *Metrics {
	lat := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "gateway_request_duration_seconds",
		Help:    "Request latency per backend",
		Buckets: prometheus.DefBuckets,
	}, []string{"backend", "method", "status"})

	req := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "gateway_requests_total",
		Help: "Total requests proxied per backend",
	}, []string{"backend", "method", "status"})

	prometheus.MustRegister(lat, req)
	return &Metrics{lat, req}
}

func (m *Metrics) Record(backend, method string, status int, d time.Duration) {
	s := strconv.Itoa(status)
	m.latency.WithLabelValues(backend, method, s).Observe(d.Seconds())
	m.requests.WithLabelValues(backend, method, s).Inc()
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}
