package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	serveMux := http.ServeMux{}
	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", apiConfig.middlewareMetricsInc(fileHandler))
	serveMux.Handle("/app/assets/logo.png", apiConfig.middlewareMetricsInc(fileHandler))
	serveMux.HandleFunc("GET /healthz", healthz)
	serveMux.HandleFunc("GET /metrics", apiConfig.getMetrics)
	serveMux.HandleFunc("POST /reset", apiConfig.resetMetrics)

	server := http.Server{
		Handler: &serveMux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(responseWriter, request)
	})
}

func (cfg *apiConfig) getMetrics(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("content-type", "  text/plain; charset=utf-8")
	responseWriter.WriteHeader(200)
	bytes := fmt.Appendf(nil, "Hits: %v", cfg.fileserverHits.Load())
	responseWriter.Write(bytes)
}

func (cfg *apiConfig) resetMetrics(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("content-type", "  text/plain; charset=utf-8")
	responseWriter.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	bytes := fmt.Appendf(nil, "Hits reset: %v", cfg.fileserverHits.Load())
	responseWriter.Write(bytes)
}
