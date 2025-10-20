package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	TC_HTTP_REQUEST_TOTAL prometheus.Counter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tarantula_http_request_total",
		Help: "Tarantual total request number",
	})
	TG_HTTP_REQUEST_DRUTAION prometheus.Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tarantula_http_request_duration",
		Help: "Tarantual request duration",
	})
)

func Register() {
	prometheus.Register(TC_HTTP_REQUEST_TOTAL)
	prometheus.Register(TG_HTTP_REQUEST_DRUTAION)
}
