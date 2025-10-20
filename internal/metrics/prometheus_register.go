package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	Counter prometheus.Counter
)

func Register() {
	Counter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Total request number",
	})
	prometheus.MustRegister(Counter)
}
