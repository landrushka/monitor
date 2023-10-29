package workers

import (
	"github.com/go-chi/chi/v5"
	"github.com/landrushka/monitor.git/internal/handlers"
	"github.com/landrushka/monitor.git/internal/logger"
	"github.com/landrushka/monitor.git/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

func StartServer(host string) error {
	var memStorage = storage.MemStorage{GaugeMetric: make(storage.GaugeMetric), CounterMetric: make(storage.CounterMetric)}
	h := handlers.NewHandler(memStorage)
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.GetAllNamesHandle)
		r.Post("/update", h.UpdateHandle)
		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", h.GetValueHandle)
		})
	})

	if err := logger.Initialize("INFO"); err != nil {
		panic("cannot initialize zap")
	}

	logger.Log.Info("Running server", zap.String("address", host))

	return http.ListenAndServe(host, r)
}
