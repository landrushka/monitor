package storage

type Storage interface {
	UpdateGauge(n string, v float64) string
	UpdateCounter(n string, v float64) string
}

type GaugeMetric map[string]float64
type CounterMetric map[string]int64

type MemStorage struct {
	GaugeMetric   GaugeMetric
	CounterMetric CounterMetric
}

func (gm MemStorage) UpdateGauge(n string, v float64) {
	_, ok := gm.GaugeMetric[n]
	if ok {
		gm.GaugeMetric[n] = v
	} else {
		gm.GaugeMetric[n] = v
	}
}

func (cm MemStorage) UpdateCounter(n string, v int64) {
	val, ok := cm.CounterMetric[n]
	if ok {
		cm.CounterMetric[n] = val + v
	} else {
		cm.CounterMetric[n] = v
	}
}
