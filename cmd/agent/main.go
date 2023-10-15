package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

const BaseURL = `http://localhost:8080`
const pollInterval = 2

//const reportInterval = 10

type statsFloat map[string]float64
type statsInt map[string]int64
type counter int64

func main() {

	// Create an HTTP client
	client := &http.Client{}
	sf := statsFloat{}
	si := statsInt{}
	count := counter(0)

	for {
		time.Sleep(pollInterval * time.Second)
		sf.getRandomMetrics()
		sf.getRandomMetrics()
		si.getCount()
		count++
		if count%5 == 0 {
			count = 0
			for name, value := range sf {
				var req, _ = makeGaugeRequest(name, value)
				resp, err := client.Do(req)
				if err != nil {
					// handle error
				}
				err = resp.Body.Close()
				if err != nil {
					return
				}
			}
			for name, value := range si {
				var req, _ = makeCounterRequest(name, value)
				resp, err := client.Do(req)
				if err != nil {
					// handle error
				}
				err = resp.Body.Close()
				if err != nil {
					return
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

func makeGaugeRequest(metricName string, metricValue float64) (*http.Request, error) {
	url := BaseURL + "/update/gauge/" + metricName + "/" + fmt.Sprintf("%.2f", metricValue)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "text/plain")
	return req, err
}

func makeCounterRequest(metricName string, metricValue int64) (*http.Request, error) {
	url := BaseURL + "/update/counter/" + metricName + "/" + strconv.FormatInt(metricValue, 10)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "text/plain")
	return req, err
}
