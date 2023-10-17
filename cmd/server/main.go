package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/landrushka/monitor.git/internal/handlers"
	"github.com/landrushka/monitor.git/internal/storage"
	"log"
	"net/http"
)

func main() {
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

	log.Fatal(http.ListenAndServe(":8080", r))
}
