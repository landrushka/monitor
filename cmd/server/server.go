package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/landrushka/monitor.git/internal/workers"
)

type Config struct {
	TargetHost string `env:"ADDRESS"`
}

var cfg Config

func main() {
	//serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	flag.StringVar(&cfg.TargetHost, "a", "localhost:8080", "Target base host:port")
	flag.Parse()

	_ = env.Parse(&cfg)

	//serverFlags.Parse(os.Args[1:])
	err := workers.StartServer(cfg.TargetHost)
	if err != nil {
		panic(err)
	}
}
