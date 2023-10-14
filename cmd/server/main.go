package main

import (
	"github.com/landrushka/monitor.git/internal/storage"
	"net/http"
	"strconv"
	"strings"

	"github.com/landrushka/monitor.git/internal/handlers"
)

var MemStorage = storage.MemStorage{GaugeMetric: make(storage.GaugeMetric), CounterMetric: make(storage.CounterMetric)}

func main() {

	mux := http.NewServeMux()
	mux.Handle(`/update/gauge/`, handlers.Middleware(http.HandlerFunc(UpdateGauge)))
	mux.Handle(`/update/counter/`, handlers.Middleware(http.HandlerFunc(UpdateCounter)))
	mux.HandleFunc(`/update/`, BadRequest)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func UpdateGauge(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	splittedPath := strings.Split(path, "/")
	// idx=2 -> type
	// idx=3 -> name
	// idx=4 -> value
	//name := splittedPath[3]
	//value := splittedPath[4]
	if len(splittedPath) <= 4 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(splittedPath[4]), 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	MemStorage.UpdateGauge(splittedPath[3], val)
	res.WriteHeader(http.StatusOK)
}

func UpdateCounter(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	splittedPath := strings.Split(path, "/")
	// idx=2 -> type
	// idx=3 -> name
	// idx=4 -> value
	//name := splittedPath[3]
	//value := splittedPath[4]
	if len(splittedPath) <= 4 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	val, err := strconv.ParseInt(strings.TrimSpace(splittedPath[4]), 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	MemStorage.UpdateCounter(splittedPath[3], val)
	res.WriteHeader(http.StatusOK)
}

func BadRequest(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}
