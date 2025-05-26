package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsPath = "/metrics"
)

var httpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"path", "method"},
)

type Middleware struct {
	enabled bool
}

func NewMiddleware(enabled bool) *Middleware {
	if enabled {
		prometheus.MustRegister(httpRequestDuration)
	}

	return &Middleware{enabled: enabled}
}

func (m *Middleware) Routes(router *http.ServeMux) {
	if !m.enabled {
		return
	}

	router.Handle(metricsPath, promhttp.Handler())
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	if !m.enabled {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.observeDuration(start, r.URL.Path, r.Method)
	})
}

func (m *Middleware) observeDuration(start time.Time, path, method string) {
	if !m.enabled {
		return
	}

	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues(path, method).Observe(duration)
}
