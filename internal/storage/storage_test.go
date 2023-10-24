package storage

import "testing"

func TestMemStorage_UpdateCounter(t *testing.T) {
	type fields struct {
		GaugeMetric   GaugeMetric
		CounterMetric CounterMetric
	}
	type args struct {
		n string
		v int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "counter test",
			fields: fields{
				GaugeMetric:   GaugeMetric{"gauge_test_name": 0.001},
				CounterMetric: CounterMetric{"counter_test_name": 1},
			},
			args: args{
				n: "counter_test_name",
				v: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := MemStorage{
				GaugeMetric:   tt.fields.GaugeMetric,
				CounterMetric: tt.fields.CounterMetric,
			}
			ms.UpdateCounter(tt.args.n, tt.args.v)
		})
	}
}

func TestMemStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		GaugeMetric   GaugeMetric
		CounterMetric CounterMetric
	}
	type args struct {
		n string
		v float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "gauge test",
			fields: fields{
				GaugeMetric:   GaugeMetric{"gauge_test_name": 0.001},
				CounterMetric: CounterMetric{"counter_test_name": 1},
			},
			args: args{
				n: "gauge_test_name",
				v: 0.001,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := MemStorage{
				GaugeMetric:   tt.fields.GaugeMetric,
				CounterMetric: tt.fields.CounterMetric,
			}
			ms.UpdateGauge(tt.args.n, tt.args.v)
			if got := ms.GaugeMetric[tt.args.n]; got != tt.args.v {
				t.Errorf("GaugeMetric = %v, want %v", got, tt.args.v)
			}
		})
	}
}
