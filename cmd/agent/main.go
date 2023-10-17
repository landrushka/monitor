package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

var targetHost = `http://localhost:8080`
var pollInterval int64 = 2
var reportInterval int64 = 10

type statsFloat map[string]float64
type statsInt map[string]int64

func main() {
	agentFlags := flag.NewFlagSet("agent", flag.ExitOnError)
	agentFlags.StringVar(&targetHost, "a", targetHost, "Target base host:port")
	agentFlags.Int64Var(&reportInterval, "r", reportInterval, "Report interval in sec")
	agentFlags.Int64Var(&pollInterval, "p", pollInterval, "Poll interval in sec")
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
		count++
		if count%reportInterval == 0 {
			count = 0
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
