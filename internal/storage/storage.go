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
	val, ok := gm.GuageMetric[n]
	if ok {
		println("exist")
		println(val)
		gm.GuageMetric[n] = v
	} else {
		println(v)
		println(n)
		gm.GuageMetric[n] = v
	}
}

func (cm MemStorage) UpdateCounter(n string, v int64) {
	val, ok := cm.CounterMetric[n]
	if ok {
		println("exist")
		println(val)
		cm.CounterMetric[n] = cm.CounterMetric[n] + v
	} else {
		println(v)
		println(n)
		cm.CounterMetric[n] = v
	}
}

func main() {

}
