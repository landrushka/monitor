package workers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	compressor := middleware.Compress(5, "text/html", "application/json")
	r.Use(handlers.GzipMiddleware, compressor, logger.RequestLogger)
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.GetAllNamesHandle)
		r.Route("/update", func(r chi.Router) {
			r.Post("/", h.UpdateHandle)
			r.Route("/{type}", func(r chi.Router) {
				r.Post("/{name}/{value}", h.UpdateHandleByParams)
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Post("/", h.GetValueHandle)
			r.Get("/{type}/{name}", h.GetValueHandleByParams)
		})
	})

	if err := logger.Initialize("INFO"); err != nil {
		panic("cannot initialize zap")
	}

	logger.Log.Info("Running server", zap.String("address", host))

	return http.ListenAndServe(host, r)
}
