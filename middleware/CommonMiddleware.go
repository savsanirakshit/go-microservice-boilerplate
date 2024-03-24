package middleware

import (
	"golang-microservice-boilerplate/logger"
	"net/http"
)

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.ServiceLogger.Debug("API call received for api : ", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
