package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	TC_HTTP_REQUEST_TOTAL prometheus.Counter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tarantula_http_request_total",
		Help: "Tarantual total request number",
	})
	TH_HTTP_REQUEST_DURATION prometheus.Histogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "tarantula_http_request_duration",
		Help:    "Tarantula http request duration",
		Buckets: prometheus.DefBuckets,
	})
	TG_SOCKET_CONCURRENT_NUMBER prometheus.Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tarantula_socket_concurrent_number",
		Help: "Tarantual socket concurrent number",
	})
)

func Register() {
	prometheus.Register(TC_HTTP_REQUEST_TOTAL)
	prometheus.Register(TG_SOCKET_CONCURRENT_NUMBER)
	prometheus.Register(TH_HTTP_REQUEST_DURATION)
}
