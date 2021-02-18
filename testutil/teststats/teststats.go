// +build testing

package teststats

import "github.com/miscoler/ethereum-proxy/pkg/stats"

type fakeStats struct{}

func New() stats.Stats {
	return &fakeStats{}
}

type fakeTiming struct{}

func (f *fakeTiming) Start() stats.Timing {
	return f
}

func (f *fakeTiming) Stop() {
}

func (f *fakeStats) NewTiming(string, string) stats.Timing {
	return &fakeTiming{}
}

type fakeMetric struct{}

func (f *fakeMetric) Observe(float64) {}

func (f *fakeStats) NewMetric(string, string) stats.Metric {
	return &fakeMetric{}
}
