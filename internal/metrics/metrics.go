package metrics

import (
	"github.com/landrushka/monitor.git/internal/storage"
	"math/rand"
	"runtime"
)

type StatsInt storage.CounterMetric
type StatsFloat storage.GaugeMetric

func (sf StatsFloat) GetRandomMetrics() {
	sf["RandomValue"] = rand.Float64()
}

func (si StatsInt) GetCount() {
	si["PollCount"] = si["PollCount"] + 1
}

func (sf StatsFloat) GetGaugeMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	sf["Alloc"] = float64(memStats.Alloc)
	sf["BuckHashSys"] = float64(memStats.BuckHashSys)
	sf["Frees"] = float64(memStats.Frees)
	sf["GCCPUFraction"] = memStats.GCCPUFraction
	sf["GCSys"] = float64(memStats.GCSys)
	sf["HeapAlloc"] = float64(memStats.HeapAlloc)
	sf["HeapIdle"] = float64(memStats.HeapIdle)
	sf["HeapInuse"] = float64(memStats.HeapInuse)
	sf["HeapObjects"] = float64(memStats.HeapObjects)
	sf["HeapReleased"] = float64(memStats.HeapReleased)
	sf["HeapSys"] = float64(memStats.HeapSys)
	sf["LastGC"] = float64(memStats.LastGC)
	sf["Lookups"] = float64(memStats.Lookups)
	sf["MCacheInuse"] = float64(memStats.MCacheInuse)
	sf["MCacheSys"] = float64(memStats.MCacheSys)
	sf["MSpanInuse"] = float64(memStats.MSpanInuse)
	sf["MSpanSys"] = float64(memStats.MSpanSys)
	sf["Mallocs"] = float64(memStats.Mallocs)
	sf["NextGC"] = float64(memStats.NextGC)
	sf["NumForcedGC"] = float64(memStats.NumForcedGC)
	sf["NumGC"] = float64(memStats.NumGC)
	sf["OtherSys"] = float64(memStats.OtherSys)
	sf["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	sf["StackInuse"] = float64(memStats.StackInuse)
	sf["StackSys"] = float64(memStats.StackSys)
	sf["Sys"] = float64(memStats.Sys)
	sf["TotalAlloc"] = float64(memStats.TotalAlloc)
}
