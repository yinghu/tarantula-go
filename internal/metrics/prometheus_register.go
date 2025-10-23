package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	BUCKET_TH_HTTP_REQUEST_DURATION                      = []float64{0.01, 0.05, 0.1, 0.5, 1, 5, 10}
	QUANTILE_TS_REQUESTION_DURATION                      = map[float64]float64{0.5: 0.02, 0.9: 0.02}
	TH_HTTP_REQUEST                 prometheus.Histogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "tarantula_http_request",
		Help:    "Tarantula http request duration",
		Buckets: BUCKET_TH_HTTP_REQUEST_DURATION,
	})
	TG_SOCKET_CONCURRENCY prometheus.Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tarantula_socket_concurrent_number",
		Help: "Tarantual socket concurrent number",
	})
	TS_HTTP_REQUEST prometheus.Summary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "tarantula_http_request_summary",
		Help:       "Tarantula http request duration summary",
		Objectives: QUANTILE_TS_REQUESTION_DURATION,
	})
	THVC_HTTP_REQUEST prometheus.HistogramVec = *prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "tarantula_http_request_vec",
		Help:    "Tarantula http request duration vec",
		Buckets: BUCKET_TH_HTTP_REQUEST_DURATION,
	}, []string{"path"})
)

func init() {
	prometheus.MustRegister(THVC_HTTP_REQUEST)
	prometheus.MustRegister(TG_SOCKET_CONCURRENCY)
	prometheus.MustRegister(TH_HTTP_REQUEST)
	prometheus.MustRegister(TS_HTTP_REQUEST)
}
