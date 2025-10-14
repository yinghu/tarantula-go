package core



type MetricsService interface {
	SetUp
	WebRequest(m ReqMetrics) error
}

type ReqMetrics struct {
	Path     string
	ReqTimed int64
	Node     string
	ReqId    int32
	ReqCode  int32
}
