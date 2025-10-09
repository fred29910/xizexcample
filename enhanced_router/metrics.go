package enhanced_router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "zinx_router_requests_total",
			Help: "Total number of requests processed by zinx router",
		},
		[]string{"path", "code"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "zinx_router_request_duration_seconds",
			Help:    "Router request processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal, requestDuration)
}

// MetricsMiddleware creates a middleware for recording Prometheus metrics.
func MetricsMiddleware() MiddlewareFunc {
	return func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			start := time.Now()
			path, _ := ctx.Value("path").(string)

			resp, err = next(ctx, req)

			duration := time.Since(start).Seconds()
			requestDuration.WithLabelValues(path).Observe(duration)

			code := "200"
			if err != nil {
				code = "500" // Or a more specific error code
			}
			requestsTotal.WithLabelValues(path, code).Inc()

			return resp, err
		}
	}
}

// StartMetricsServer starts an HTTP server to expose Prometheus metrics.
func StartMetricsServer(addr, path string) {
	http.Handle(path, promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Println("Metrics server error: ", err)
		}
	}()
}
