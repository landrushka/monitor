package handlers

import (
	"github.com/go-chi/chi/v5"
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
	typeName := strings.ToLower(chi.URLParam(r, "type"))
	if typeName != "gauge" && typeName != "counter" {
		http.Error(rw, "unknown type: "+typeName, http.StatusBadRequest)
	}
	if typeName == "gauge" {
		name := chi.URLParam(r, "name")
		value := strings.ToLower(chi.URLParam(r, "value"))
		val, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		h.memStorage.UpdateGauge(name, val)
		rw.WriteHeader(http.StatusOK)
	}
	if typeName == "counter" {
		name := chi.URLParam(r, "name")
		value := strings.ToLower(chi.URLParam(r, "value"))
		val, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		h.memStorage.UpdateCounter(name, val)
		rw.WriteHeader(http.StatusOK)
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
