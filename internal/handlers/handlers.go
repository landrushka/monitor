package handlers

import (
	"net/http"
)

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
