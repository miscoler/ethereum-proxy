package stats

import (
	"net/http"
	"runtime"
	"time"

	"github.com/miscoler/ethereum-proxy/pkg/logger"
	"github.com/miscoler/ethereum-proxy/pkg/stats/config"

	"github.com/benbjohnson/clock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/load"
	"go.uber.org/zap"
)

type Stats interface {
	NewMetric(string, string) Metric
	NewTiming(string, string) Timing
}

type Metric interface {
	Observe(float64)
}

type Timing interface {
	Start() Timing
	Stop()
}

type statsImpl struct {
	clock      clock.Clock
	rss        Metric
	la         Metric
	goroutines Metric
}

type metricImpl struct {
	prometheus.Histogram
}

func (m *metricImpl) Observe(v float64) {
	m.Histogram.Observe(v)
}

func (s *statsImpl) NewMetric(name string, desc string) Metric {
	return &metricImpl{
		Histogram: promauto.NewHistogram(prometheus.HistogramOpts{ //nolint: exhaustivestruct
			Name: name,
			Help: desc,
		}),
	}
}

type timingImpl struct {
	clock clock.Clock
	prometheus.Summary
	startTime time.Time
}

func (t *timingImpl) Start() Timing {
	t.startTime = t.clock.Now()

	return t
}

func (t *timingImpl) Stop() {
	t.Observe(float64(t.clock.Since(t.startTime)))
}

func (s *statsImpl) NewTiming(name string, desc string) Timing {
	return &timingImpl{
		clock: s.clock,
		Summary: promauto.NewSummary(prometheus.SummaryOpts{ //nolint: exhaustivestruct
			Name: name,
			Help: desc,
		}),
	}
}

func New(
	conf config.Config,
	clock clock.Clock,
	lg logger.Logger,
) (Stats, error) {
	st := &statsImpl{
		clock: clock,
	}
	st.rss = st.NewMetric(
		"rss",
		"rss of current process",
	)
	st.la = st.NewMetric(
		"la",
		"load average 1min of current process",
	)
	st.goroutines = st.NewMetric(
		"goroutines",
		"number of goroutines in current process",
	)

	go func() {
		for range time.Tick(time.Minute) {
			lg.Info("sending runtime stats")
			st.goroutines.Observe(float64(runtime.NumGoroutine()))
			st.la.Observe(float64(runtime.NumCPU()))
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			st.rss.Observe(float64(m.Alloc))
			load, err := load.Avg()
			if err != nil {
				lg.Error("error getting load avg", zap.Error(err))

				continue
			}
			st.la.Observe(load.Load1)
		}
	}()

	go func() {
		http.Handle(conf.Endpoint, promhttp.Handler())
		if err := http.ListenAndServe(conf.Addr, nil); err != nil {
			lg.Error("error starting metrics server", zap.Error(err))
		}
	}()

	return st, nil
}
