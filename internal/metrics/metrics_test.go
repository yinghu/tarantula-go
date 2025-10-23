package metrics

import (
	"fmt"
	"strings"
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
	bkt := []float64{0.1, 0.5, 1.0}
	hg := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "test_histogram",
		Help:    "Test histogram",
		Buckets: bkt,
	})
	prometheus.Register(hg)
	hg.Observe(0.05)
	hg.Observe(0.2)
	hg.Observe(0.7)
	hg.Observe(1.5)
	expected := `
			# HELP test_histogram Test histogram
    		# TYPE test_histogram histogram
    		test_histogram_bucket{le="0.1"} 1
    		test_histogram_bucket{le="0.5"} 2
    		test_histogram_bucket{le="1.0"} 3
    		test_histogram_bucket{le="+Inf"} 4 
    		test_histogram_sum 2.45
    		test_histogram_count 4
    	`
	err := testutil.CollectAndCompare(hg, strings.NewReader(expected))
	if err != nil {
		fmt.Printf("Error %s\n", err.Error())
	}
	var vec *prometheus.HistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "",
		Help:    "",
		Buckets: bkt,
	}, []string{"path"})
	vec.WithLabelValues("/abc").Observe(10)
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
