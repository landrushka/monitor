package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatsFloat_GetGaugeMetrics(t *testing.T) {
	tests := []struct {
		name string
		sf   StatsFloat
	}{
		{name: "positive test #1",
			sf: StatsFloat{"gauge_test_name": 0.001}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sf.GetGaugeMetrics()
			assert.IsType(t, float64(1), tt.sf["Alloc"])
		})
	}
}

func TestStatsFloat_GetRandomMetrics(t *testing.T) {
	tests := []struct {
		name string
		sf   StatsFloat
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sf.GetRandomMetrics()
		})
	}
}

func TestStatsInt_GetCount(t *testing.T) {
	tests := []struct {
		name string
		si   StatsInt
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.si.GetCount()
		})
	}
}
