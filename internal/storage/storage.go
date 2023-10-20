package storage

type Storage interface {
	UpdateGauge(n string, v float64)
	UpdateCounter(n string, v int64)
}

type GaugeMetric map[string]float64
type CounterMetric map[string]int64

type StatsFloat map[string]float64
type StatsInt map[string]int64

type MemStorage struct {
	GaugeMetric   GaugeMetric
	CounterMetric CounterMetric
}

func (ms MemStorage) UpdateGauge(n string, v float64) {
	_, ok := ms.GaugeMetric[n]
	if ok {
		ms.GaugeMetric[n] = v
	} else {
		ms.GaugeMetric[n] = v
	}
}

func (ms MemStorage) UpdateCounter(n string, v int64) {
	val, ok := ms.CounterMetric[n]
	if ok {
		ms.CounterMetric[n] = val + v
	} else {
		ms.CounterMetric[n] = v
	}
}
