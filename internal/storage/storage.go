package storage

import (
	"encoding/json"
	"github.com/landrushka/monitor.git/internal/logger"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"runtime"
)

type Storage interface {
	UpdateGauge(n string, v float64)
	UpdateCounter(n string, v int64)
	SaveData()
	RestoreData()
}

type GaugeMetric map[string]float64
type CounterMetric map[string]int64

type MemStorage struct {
	GaugeMetric   GaugeMetric
	CounterMetric CounterMetric
	Producer      *Producer
	Consumer      *Consumer
}

func NewMemStorage() (*MemStorage, error) {
	return &MemStorage{GaugeMetric: make(GaugeMetric), CounterMetric: make(CounterMetric)}, nil
}

var met Metrics

func (ms MemStorage) SaveData() {
	for k, v := range ms.GaugeMetric {
		met.ID = k
		met.MType = "gauge"
		met.Value = &v
		logger.Log.Info("SaveData", zap.String("name", k), zap.Float64("value", v))
		_ = ms.Producer.WriteMetric(&met)
	}
	for k, v := range ms.CounterMetric {
		met.ID = k
		met.MType = "counter"
		met.Delta = &v
		logger.Log.Info("SaveData", zap.String("name", k), zap.Int64("value", v))
		_ = ms.Producer.WriteMetric(&met)
	}
	//_ = ms.producer.Close()
}

func (ms MemStorage) RestoreData() {
	for {
		m, err := ms.Consumer.ReadMetric()
		if err != nil {
			_ = ms.Consumer.Close()
			return
		}
		switch typeName := m.MType; typeName {
		case "gauge":
			ms.UpdateGauge(m.ID, *m.Value)
		case "counter":
			ms.UpdateCounter(m.ID, *m.Delta)
		default:
			return
		}
	}
	_ = ms.Consumer.Close()
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

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (sf GaugeMetric) GetRandomMetrics() {
	sf["RandomValue"] = rand.Float64()
}

func (si CounterMetric) GetCount() {
	si["PollCount"] = si["PollCount"] + 1
}

func (sf GaugeMetric) GetGaugeMetrics() {
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

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteMetric(metric *Metrics) error {
	return p.encoder.Encode(&metric)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadMetric() (*Metrics, error) {
	metric := &Metrics{}
	if err := c.decoder.Decode(&metric); err != nil {
		return nil, err
	}

	return metric, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
