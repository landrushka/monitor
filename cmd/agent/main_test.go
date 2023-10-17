package agent

import "testing"

func Test_statsFloat_getGaugeMetrics(t *testing.T) {
	tests := []struct {
		name string
		sf   statsFloat
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sf.getGaugeMetrics()
		})
	}
}
