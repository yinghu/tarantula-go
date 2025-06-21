package metrics

import (
	"gameclustering.com/internal/core"
)

type MetricsService interface {
	core.SetUp
	WebRequest(m ReqMetrics) error
}

type ReqMetrics struct {
	Path     string
	ReqTimed int64
	Node     string
	ReqId    int32
	ReqCode  int32
}
