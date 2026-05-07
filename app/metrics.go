package main

import (
	"net/http"
	"fmt"
)

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	val := fmt.Sprintf(adminstring, cfg.fileserverHits.Load())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(val))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}