package metrics

import (
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestCouner(t *testing.T) {
	ct := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "Test counter",
	})
	ct.Inc()
	ct.Add(101)
	n := testutil.ToFloat64(ct)
	if n != 102 {
		t.Errorf("should be 102 %.f\n", n)
	}

	vct := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "Test counter",
	}, []string{"path", "method"})
	vct.WithLabelValues("a", "c").Add(1)
	vct.WithLabelValues("a", "c").Add(1)
	vct.WithLabelValues("a", "c").Add(1)
	vct.WithLabelValues("b", "c").Add(1)
	vct.WithLabelValues("b", "c").Add(1)
	vct.WithLabelValues("b", "c").Add(1)
	fmt.Printf("Total %.f\n", testutil.ToFloat64(vct.WithLabelValues("a", "c")))
	fmt.Printf("Total %.f\n", testutil.ToFloat64(vct.WithLabelValues("b", "c")))
}

func TestGauge(t *testing.T) {
	ct := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "test_counter",
		Help: "Test counter",
	})
	ct.Inc()
	ct.Add(101)
	ct.Dec()
	n := testutil.ToFloat64(ct)
	if n != 101 {
		t.Errorf("should be 101 %.f\n", n)
	}

	vct := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test_gauge",
		Help: "Test gauge",
	}, []string{"path", "method"})
	vct.WithLabelValues("a", "c").Add(1)
	vct.WithLabelValues("a", "c").Add(1)
	vct.WithLabelValues("a", "c").Add(1)
	vct.WithLabelValues("b", "c").Add(1)
	vct.WithLabelValues("b", "c").Add(1)
	vct.WithLabelValues("b", "c").Add(1)
	vct.WithLabelValues("b", "c").Sub(3)
	vct.WithLabelValues("a", "c").Sub(3)
	fmt.Printf("Total %.f\n", testutil.ToFloat64(vct.WithLabelValues("a", "c")))
	fmt.Printf("Total %.f\n", testutil.ToFloat64(vct.WithLabelValues("b", "c")))
}

func TestHistogram(t *testing.T) {
	bkt := []float64{50, 100, 200, 300, 400, 500}
	hg := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "test_histogram",
		Help:    "Test histogram",
		Buckets: bkt,
	})
	hg.Observe(306)
	hg.Observe(306)
	ct := testutil.CollectAndCount(hg, `test_histogram_bucket{le="500"}`)
	fmt.Printf("Error %d\n", ct)

}

func TestSummay(t *testing.T) {
	pct := map[float64]float64{0.5: 0.02, 0.9: 0.01}
	hg := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "test_summary",
		Help:       "Test summary",
		Objectives: pct,
	})
	hg.Observe(200)
	hg.Observe(160)
	ct := testutil.CollectAndCount(hg, `test_histogram_bucket{le="500"}`)
	fmt.Printf("Error %d\n", ct)

}
