package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/landrushka/monitor.git/internal/handlers"
	"github.com/landrushka/monitor.git/internal/storage"
	"log"
	"net/http"
	"os"
)

var targetHost string

type Config struct {
	targetHost string `env:"ADDRESS" envDefault:":8080"`
}

func main() {
	//serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	var cfg Config
	_ = env.Parse(&cfg)
	targetHost = os.Getenv("ADDRESS")
	flag.StringVar(&targetHost, "a", cfg.targetHost, "Target base host:port")
	flag.Parse()
	//serverFlags.Parse(os.Args[1:])
	var memStorage = storage.MemStorage{GaugeMetric: make(storage.GaugeMetric), CounterMetric: make(storage.CounterMetric)}
	h := handlers.NewHandler(memStorage)

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.GetAllNamesHandle)
		r.Route("/update", func(r chi.Router) {
			r.Route("/{type}", func(r chi.Router) {
				r.Post("/{name}/{value}", h.UpdateHandle)
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", h.GetValueHandle)
		})
	})

	log.Fatal(http.ListenAndServe(targetHost, r))
}
