package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/landrushka/monitor.git/internal/handlers"
	"github.com/landrushka/monitor.git/internal/storage"
	"log"
	"net/http"
)

var targetHost string

type Config struct {
	TargetHost string `env:"ADDRES" envDefault:":8080"`
}

func main() {
	//serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	cfg := Config{}
	_ = env.Parse(&cfg)

	flag.StringVar(&targetHost, "a", cfg.TargetHost, "Target base host:port")
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
