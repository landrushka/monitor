package workers

import "testing"

func TestStartWorker(t *testing.T) {
	type args struct {
		host           string
		reportInterval int64
		pollInterval   int64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartWorker(tt.args.host, tt.args.reportInterval, tt.args.pollInterval)
		})
	}
}
