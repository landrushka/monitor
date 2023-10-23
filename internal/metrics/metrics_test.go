package metrics

import "testing"

func TestStatsFloat_GetGaugeMetrics(t *testing.T) {
	tests := []struct {
		name string
		sf   StatsFloat
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sf.GetGaugeMetrics()
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
