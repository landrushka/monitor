package main

import (
	"github.com/landrushka/monitor.git/internal/storage"
	"net/http"
	"strconv"
	"strings"

	"github.com/landrushka/monitor.git/internal/handlers"
)

var MemStorage = storage.MemStorage{GuageMetric: make(storage.GuageMetric), CounterMetric: make(storage.CounterMetric)}

func main() {

	mux := http.NewServeMux()
	mux.Handle(`/update/guage/`, handlers.Middleware(http.HandlerFunc(UpdateGuage)))
	mux.Handle(`/update/counter/"`, handlers.Middleware(http.HandlerFunc(UpdateCounter)))
	mux.HandleFunc(`/update/`, BadRequest)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func UpdateGuage(res http.ResponseWriter, req *http.Request) {
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
	val, _ := strconv.ParseFloat(strings.TrimSpace(splittedPath[4]), 64)
	MemStorage.UpdateGuage(splittedPath[3], val)
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
	val, _ := strconv.ParseInt(strings.TrimSpace(splittedPath[4]), 64, 64)
	MemStorage.UpdateCounter(splittedPath[3], val)
	res.WriteHeader(http.StatusOK)
}

func BadRequest(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}
