package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/landrushka/monitor.git/internal/workers"
	"sync"
	"time"
)

type Config struct {
	TargetHost      string        `env:"ADDRESS"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
}

var cfg Config

func main() {
	flag.StringVar(&cfg.TargetHost, "a", "localhost:8080", "Target base host:port")
	flag.DurationVar(&cfg.StoreInterval, "i", 300, "Report interval in sec")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "FileStoragePath")
	flag.BoolVar(&cfg.Restore, "r", true, "Restore")
	flag.Parse()
	_ = env.Parse(&cfg)

	uptimeTicker := time.NewTicker(cfg.StoreInterval * time.Second)
	var saveNow = false
	if cfg.StoreInterval == 0 {
		saveNow = true
		uptimeTicker = nil
	}
	workers.InitFileManager(cfg.Restore, cfg.FileStoragePath)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := workers.StartServer(cfg.TargetHost, saveNow)
		if err != nil {
			workers.StartFileManager()
			panic(err)
		}
	}()
	for {
		select {
		case <-uptimeTicker.C:
			go func() {
				workers.StartFileManager()
			}()
		}
	}
	//wg.Wait()

}
