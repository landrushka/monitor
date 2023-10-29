package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/landrushka/monitor.git/internal/metrics"
	"github.com/landrushka/monitor.git/internal/storage"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const nameListHTML = `
<h1>Metric Names</h1>
<dl>
{{.}}
</dl>
`

func NewHandler(memStorage storage.MemStorage) *Handler {
	h := &Handler{
		memStorage: memStorage,
	}
	return h
}

type Handler struct {
	memStorage storage.MemStorage
}

func (h *Handler) UpdateHandle(rw http.ResponseWriter, r *http.Request) {
	var m metrics.Metrics
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	switch typeName := m.MType; typeName {
	case "gauge":
		h.memStorage.UpdateGauge(m.ID, *m.Value)
		json.NewEncoder(rw).Encode(m)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	case "counter":
		h.memStorage.UpdateCounter(m.ID, *m.Delta)
		json.NewEncoder(rw).Encode(m)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	default:
		http.Error(rw, "unknown type: "+typeName, http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetAllNamesHandle(rw http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(h.memStorage.GaugeMetric)+len(h.memStorage.CounterMetric))
	for k := range h.memStorage.GaugeMetric {
		keys = append(keys, k)
	}
	for k := range h.memStorage.CounterMetric {
		keys = append(keys, k)
	}
	tmpl := template.Must(template.New("").Parse(nameListHTML))
	_ = tmpl.Execute(rw, keys)
}

func (h *Handler) GetValueHandle(rw http.ResponseWriter, r *http.Request) {
	typeName := strings.ToLower(chi.URLParam(r, "type"))
	nameName := chi.URLParam(r, "name")
	if typeName == "gauge" {
		val, ok := h.memStorage.GaugeMetric[nameName]
		if ok {
			_, _ = rw.Write([]byte(strconv.FormatFloat(val, 'f', -1, 64)))
		} else {
			http.Error(rw, "unknown name: "+nameName, http.StatusNotFound)
		}

	}
	if typeName == "counter" {
		val, ok := h.memStorage.CounterMetric[nameName]
		if ok {
			_, _ = rw.Write([]byte(strconv.FormatInt(val, 10)))
		} else {
			http.Error(rw, "unknown name: "+nameName, http.StatusNotFound)
		}

	}
}
