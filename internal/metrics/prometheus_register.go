package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	BUCKET_HTTP_REQUEST_DURATION = []float64{0.01, 0.05, 0.1, 0.5, 1, 5, 10}

	SOCKET_CONCURRENCY_METRICS prometheus.Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tarantula_socket_concurrent_number",
		Help: "Tarantual socket concurrent number",
	})
	HTTP_REQUEST_METRICS prometheus.HistogramVec = *prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "tarantula_http_request",
		Help:    "Tarantula http request duration",
		Buckets: BUCKET_HTTP_REQUEST_DURATION,
	}, []string{"path"})
)

func init() {
	prometheus.MustRegister(SOCKET_CONCURRENCY_METRICS)
	prometheus.MustRegister(HTTP_REQUEST_METRICS)
}
