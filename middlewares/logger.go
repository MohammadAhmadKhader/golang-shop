package middlewares

import (
	"log"
	"net/http"
	"time"

	"main.go/types"
)

func Logger(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		AppResponse := &types.AppResponse{ResponseWriter: w}

		next.ServeHTTP(AppResponse, r)

		duration := time.Since(start)
		log.Printf("%s %s %v %v", r.Method, r.RequestURI, duration, AppResponse.StatusCode)
	})
}