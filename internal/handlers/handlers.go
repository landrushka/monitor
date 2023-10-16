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

var MemStorage = storage.MemStorage{GaugeMetric: make(storage.GaugeMetric), CounterMetric: make(storage.CounterMetric)}

func UpdateHandle(rw http.ResponseWriter, r *http.Request) {
	typeName := strings.ToLower(chi.URLParam(r, "type"))
	if typeName != "gauge" && typeName != "counter" {
		http.Error(rw, "unknown type: "+typeName, http.StatusBadRequest)
	}
	if typeName == "gauge" {
		name := strings.ToLower(chi.URLParam(r, "name"))
		value := strings.ToLower(chi.URLParam(r, "value"))
		val, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		MemStorage.UpdateGauge(name, val)
		rw.WriteHeader(http.StatusOK)
	}
	if typeName == "counter" {
		name := strings.ToLower(chi.URLParam(r, "name"))
		value := strings.ToLower(chi.URLParam(r, "value"))
		val, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		MemStorage.UpdateCounter(name, val)
		rw.WriteHeader(http.StatusOK)
	}
}

func GetAllNamesHandle(rw http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(MemStorage.GaugeMetric)+len(MemStorage.CounterMetric))
	for k := range MemStorage.GaugeMetric {
		keys = append(keys, k)
	}
	for k := range MemStorage.CounterMetric {
		keys = append(keys, k)
	}
	tmpl := template.Must(template.New("").Parse(nameListHTML))
	_ = tmpl.Execute(rw, keys)
}

func GetValueHandle(rw http.ResponseWriter, r *http.Request) {
	typeName := strings.ToLower(chi.URLParam(r, "type"))
	nameName := strings.ToLower(chi.URLParam(r, "name"))
	if typeName == "gauge" {
		val, err := MemStorage.GaugeMetric[nameName]
		if err {
			http.Error(rw, "unknown name: "+typeName, http.StatusNotFound)
		}
		_, _ = rw.Write(Float64ToByte(val))
		rw.WriteHeader(http.StatusOK)
	}
	if typeName == "counter" {
		val, err := MemStorage.CounterMetric[nameName]
		if err {
			http.Error(rw, "unknown name: "+typeName, http.StatusNotFound)
		}
		_, _ = rw.Write([]byte(strconv.FormatInt(val, 10)))
		rw.WriteHeader(http.StatusOK)
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
