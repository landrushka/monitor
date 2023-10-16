package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/landrushka/monitor.git/internal/handlers"
	"log"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAllNamesHandle)
		r.Route("/update", func(r chi.Router) {
			r.Route("/{type}", func(r chi.Router) {
				r.Post("/{name}/{value}", handlers.UpdateHandle)
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", handlers.GetValueHandle)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
