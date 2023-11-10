package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/landrushka/monitor.git/internal/logger"
	"github.com/landrushka/monitor.git/internal/storage"
	"github.com/landrushka/monitor.git/internal/utils"
	"go.uber.org/zap"
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

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		//// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		//acceptEncoding := r.Header.Get("Accept-Encoding")
		//supportsGzip := strings.Contains(acceptEncoding, "gzip")
		//contentType := r.Header.Get("Content-Type")
		//supportsContentType := false
		//if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") {
		//	supportsContentType = true
		//}
		//if supportsGzip && supportsContentType {
		//
		//	// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
		//	cw := utils.NewCompressWriter(w)
		//	// меняем оригинальный http.ResponseWriter на новый
		//	ow = cw
		//	// не забываем отправить клиенту все сжатые данные после завершения middleware
		//	defer cw.Close()
		//}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := utils.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		h.ServeHTTP(ow, r)
	})
}

func NewHandler(memStorage storage.MemStorage, saveNow bool) *Handler {
	h := &Handler{
		memStorage: memStorage,
		saveNow:    saveNow,
	}
	return h
}

type Handler struct {
	memStorage storage.MemStorage
	saveNow    bool
}

func (h *Handler) UpdateHandleByParams(rw http.ResponseWriter, r *http.Request) {
	switch typeName := strings.ToLower(chi.URLParam(r, "type")); typeName {
	case "gauge":
		name := chi.URLParam(r, "name")
		value := strings.ToLower(chi.URLParam(r, "value"))
		val, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		h.memStorage.UpdateGauge(name, val)
		if h.saveNow {
			h.memStorage.SaveData()
		}
	case "counter":
		name := chi.URLParam(r, "name")
		value := strings.ToLower(chi.URLParam(r, "value"))
		val, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		h.memStorage.UpdateCounter(name, val)
		if h.saveNow {
			h.memStorage.SaveData()
		}
	default:
		http.Error(rw, "unknown type: "+typeName, http.StatusBadRequest)
	}
}

func (h *Handler) UpdateHandle(rw http.ResponseWriter, r *http.Request) {
	var m storage.Metrics
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	switch typeName := m.MType; typeName {
	case "gauge":
		if m.Value == nil {
			http.Error(rw, "", http.StatusBadRequest)
			return
		}
		h.memStorage.UpdateGauge(m.ID, *m.Value)
		if err := json.NewEncoder(rw).Encode(m); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		} else {
			if f, ok := rw.(http.Flusher); ok {
				f.Flush()
			}
		}
		if h.saveNow {
			h.memStorage.SaveData()
		}
	case "counter":
		if m.Delta == nil {
			http.Error(rw, "", http.StatusBadRequest)
			return
		}
		h.memStorage.UpdateCounter(m.ID, *m.Delta)
		if err := json.NewEncoder(rw).Encode(m); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		if h.saveNow {
			h.memStorage.SaveData()
		}
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
	rw.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.New("").Parse(nameListHTML))
	_ = tmpl.Execute(rw, keys)
}

func (h *Handler) GetValueHandle(rw http.ResponseWriter, r *http.Request) {
	var m storage.Metrics
	logger.Log.Info("Request",
		zap.String("URI", r.RequestURI),
		zap.String("method", r.Method),
		zap.String("request_Accept-Encoding", r.Header.Get("Accept-Encoding")),
		zap.String("request_Content-Encoding", r.Header.Get("Content-Encoding")),
	)
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	var n = m.ID
	switch typeName := m.MType; typeName {
	case "gauge":
		val, ok := h.memStorage.GaugeMetric[n]
		if ok {
			m.Value = &val
			if err := json.NewEncoder(rw).Encode(m); err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			} else {
				rw.WriteHeader(http.StatusOK)
			}
		} else {
			http.Error(rw, "unknown name: "+n, http.StatusNotFound)
		}
	case "counter":
		val, ok := h.memStorage.CounterMetric[n]
		if ok {
			m.Delta = &val
			if err := json.NewEncoder(rw).Encode(m); err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			} else {
				rw.WriteHeader(http.StatusOK)
			}
		} else {
			http.Error(rw, "unknown name: "+n, http.StatusNotFound)
		}
	default:
		http.Error(rw, "unknown type: "+typeName, http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetValueHandleByParams(rw http.ResponseWriter, r *http.Request) {
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
