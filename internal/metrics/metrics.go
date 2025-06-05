package metrics

type MetricsService interface {
	WebRequest(m ReqMetrics) error
}

type ReqMetrics struct {
	Path     string
	ReqTimed int64
	Node     string
	ReqId    int32
}
