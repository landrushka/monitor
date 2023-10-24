package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/landrushka/monitor.git/internal/workers"
)

type Config struct {
	TargetHost     string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
}

var cfg Config

func main() {

	flag.StringVar(&cfg.TargetHost, "a", "http://localhost:8080", "Target base host:port")
	flag.Int64Var(&cfg.ReportInterval, "r", 2, "Report interval in sec")
	flag.Int64Var(&cfg.PollInterval, "p", 10, "Poll interval in sec")
	flag.Parse()
	_ = env.Parse(&cfg)

	workers.StartWorker(cfg.TargetHost, cfg.ReportInterval, cfg.PollInterval)
}
