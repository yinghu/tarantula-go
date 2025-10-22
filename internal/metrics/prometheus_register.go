package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	BUCKET_TH_HTTP_REQUEST_DURATION = []float64{50, 100, 200, 300, 400, 500}

	TH_HTTP_REQUEST prometheus.Histogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "tarantula_http_request",
		Help:    "Tarantula http request duration",
		Buckets: BUCKET_TH_HTTP_REQUEST_DURATION,
	})
	TG_SOCKET_CONCURRENCY prometheus.Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tarantula_socket_concurrent_number",
		Help: "Tarantual socket concurrent number",
	})
)

func init() {
	//prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	prometheus.Register(TG_SOCKET_CONCURRENCY)
	prometheus.Register(TH_HTTP_REQUEST)
}
