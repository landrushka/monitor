package storage

type Storage interface {
	UpdateGuage(n string, v float64) string
	UpdateCounter(n string, v float64) string
}

type GuageMetric map[string]float64
type CounterMetric map[string]int64

type MemStorage struct {
	GuageMetric   GuageMetric
	CounterMetric CounterMetric
}

func (gm MemStorage) UpdateGuage(n string, v float64) {
	_, ok := gm.GuageMetric[n]
	if ok {
		gm.GuageMetric[n] = v
	} else {
		gm.GuageMetric[n] = v
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
