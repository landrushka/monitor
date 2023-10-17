package handlers

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

// middleware принимает параметром Handler и возвращает тоже Handler.
func Middleware(next http.Handler) http.Handler {
	// получаем Handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// здесь пишем логику обработки
		// например, разрешаем запросы cross-domain
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// ...
		// замыкание: используем ServeHTTP следующего хендлера
		if req.Method != http.MethodPost {
			http.Error(res, "Only POST requests", http.StatusMethodNotAllowed)
			return
		} else {
			next.ServeHTTP(res, req)
		}

	})
}

func NewHandler(memStorage storage.MemStorage) *Handler {
	h := &Handler{
		memStorage: memStorage,
	}
	return h
}

type Handler struct {
	memStorage storage.MemStorage
}

func (bh *Handler) UpdateHandle(rw http.ResponseWriter, r *http.Request) {
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
		bh.memStorage.UpdateGauge(name, val)
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
		bh.memStorage.UpdateCounter(name, val)
		rw.WriteHeader(http.StatusOK)
	}
}

func (bh *Handler) GetAllNamesHandle(rw http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(bh.memStorage.GaugeMetric)+len(bh.memStorage.CounterMetric))
	for k := range bh.memStorage.GaugeMetric {
		keys = append(keys, k)
	}
	for k := range bh.memStorage.CounterMetric {
		keys = append(keys, k)
	}
	tmpl := template.Must(template.New("").Parse(nameListHTML))
	_ = tmpl.Execute(rw, keys)
}

func (bh *Handler) GetValueHandle(rw http.ResponseWriter, r *http.Request) {
	typeName := strings.ToLower(chi.URLParam(r, "type"))
	nameName := chi.URLParam(r, "name")
	if typeName == "gauge" {
		val, ok := bh.memStorage.GaugeMetric[nameName]
		if ok {
			_, _ = rw.Write(Float64ToByte(val))
		} else {
			http.Error(rw, "unknown name: "+typeName, http.StatusNotFound)
		}

	}
	if typeName == "counter" {
		val, ok := bh.memStorage.CounterMetric[nameName]
		if ok {
			_, _ = rw.Write([]byte(strconv.FormatInt(val, 10)))
		} else {
			http.Error(rw, "unknown name: "+typeName, http.StatusNotFound)
		}

	}
}

func Float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}
