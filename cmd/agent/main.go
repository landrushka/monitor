package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

var targetHost string
var pollInterval int64
var reportInterval int64

type statsFloat map[string]float64
type statsInt map[string]int64

type Config struct {
	targetHost     string `env:"ADDRESS" envDefault:"http://localhost:8080"`
	reportInterval int64  `env:"REPORT_INTERVAL" envDefault:"2"`
	pollInterval   int64  `env:"POLL_INTERVAL" envDefault:"10"`
}

func main() {
	var cfg Config
	_ = env.Parse(&cfg)
	//agentFlags := flag.NewFlagSet("agent", flag.ExitOnError)
	flag.StringVar(&targetHost, "a", cfg.targetHost, "Target base host:port")
	flag.Int64Var(&reportInterval, "r", cfg.reportInterval, "Report interval in sec")
	flag.Int64Var(&pollInterval, "p", cfg.pollInterval, "Poll interval in sec")
	flag.Parse()
	client := resty.New()
	sf := statsFloat{}
	si := statsInt{}
	count := int64(0)

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		sf.getRandomMetrics()
		sf.getGaugeMetrics()
		si.getCount()
		count += pollInterval
		if count/reportInterval >= 1 {
			count -= reportInterval
			for name, value := range sf {
				_, err := client.R().SetPathParams(map[string]string{
					"name":  name,
					"value": fmt.Sprintf("%.2f", value),
				}).SetHeader("Content-Type", "text/plain").
					Post(targetHost + "/gauge/{name}/{value}")
				if err != nil {
					panic(err)
				}
			}
			for name, value := range si {
				_, err := client.R().SetPathParams(map[string]string{
					"name":  name,
					"value": strconv.FormatInt(value, 10),
				}).SetHeader("Content-Type", "text/plain").
					Post(targetHost + "/counter/{name}/{value}")
				if err != nil {
					panic(err)
				}
			}
		}

	}
}

func (sf statsFloat) getRandomMetrics() {
	sf["RandomValue"] = rand.Float64()
}

func (si statsInt) getCount() {
	si["PollCount"] = si["PollCount"] + 1
}

func (sf statsFloat) getGaugeMetrics() {
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
