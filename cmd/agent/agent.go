package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	"github.com/landrushka/monitor.git/internal/metrics"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	TargetHost     string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
}

var cfg Config

func main() {

	//agentFlags := flag.NewFlagSet("agent", flag.ExitOnError)
	flag.StringVar(&cfg.TargetHost, "a", "http://localhost:8080", "Target base host:port")
	flag.Int64Var(&cfg.ReportInterval, "r", 2, "Report interval in sec")
	flag.Int64Var(&cfg.PollInterval, "p", 10, "Poll interval in sec")
	flag.Parse()

	_ = env.Parse(&cfg)

	client := resty.New()

	sf := metrics.StatsFloat{}
	si := metrics.StatsInt{}
	count := int64(0)

	for {
		time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
		sf.GetRandomMetrics()
		sf.GetGaugeMetrics()
		si.GetCount()
		count += cfg.PollInterval
		if count/cfg.ReportInterval >= 1 {
			count -= cfg.ReportInterval
			if !strings.Contains(cfg.TargetHost, "http://") {
				cfg.TargetHost = "http://" + cfg.TargetHost
			}
			for name, value := range sf {
				_, err := client.R().SetPathParams(map[string]string{
					"name":  name,
					"value": fmt.Sprintf("%.2f", value),
				}).SetHeader("Content-Type", "text/plain").
					Post(cfg.TargetHost + "/update/gauge/{name}/{value}")
				if err != nil {
					panic(err)
				}
			}
			for name, value := range si {
				if !strings.Contains(cfg.TargetHost, "http://") {
					cfg.TargetHost = "http://" + cfg.TargetHost
				}
				_, err := client.R().SetPathParams(map[string]string{
					"name":  name,
					"value": strconv.FormatInt(value, 10),
				}).SetHeader("Content-Type", "text/plain").
					Post(cfg.TargetHost + "/update/counter/{name}/{value}")
				if err != nil {
					panic(err)
				}
			}
		}

	}
}
