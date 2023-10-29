package workers

import "testing"

func TestStartAgent(t *testing.T) {
	type args struct {
		host           string
		reportInterval int64
		pollInterval   int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StartAgent(tt.args.host, tt.args.reportInterval, tt.args.pollInterval); (err != nil) != tt.wantErr {
				t.Errorf("StartAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
