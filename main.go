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
	serveMux.HandleFunc("GET /api/healthz", healthz)
	//serveMux.HandleFunc("POST /api/validate_chirp", chirp)
	serveMux.HandleFunc("GET /admin/metrics", apiConfig.getMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiConfig.resetMetrics)
	serveMux.HandleFunc("POST /api/validate_chirp", validate_chirp)

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
	responseWriter.Header().Set("content-type", "  text/html")
	responseWriter.WriteHeader(200)
	html := fmt.Sprintf(`
	<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, cfg.fileserverHits.Load())
	bytes := fmt.Appendf(nil, "%s", html)
	responseWriter.Write(bytes)
}

func (cfg *apiConfig) resetMetrics(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("content-type", "  text/plain; charset=utf-8")
	responseWriter.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	bytes := fmt.Appendf(nil, "Hits reset: %v", cfg.fileserverHits.Load())
	responseWriter.Write(bytes)
}
