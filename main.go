package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/SebastianGeroli/chirpy-boot-dev/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Failed to open connection to DB")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}
	serveMux := http.ServeMux{}
	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))
	serveMux.Handle("/app/assets/logo.png", cfg.middlewareMetricsInc(fileHandler))
	serveMux.HandleFunc("GET /api/healthz", healthz)
	serveMux.HandleFunc("GET /admin/metrics", cfg.getMetrics)
	serveMux.HandleFunc("POST /admin/reset", cfg.reset)
	serveMux.HandleFunc("POST /api/validate_chirp", validate_chirp)
	serveMux.HandleFunc("POST /api/users", cfg.createUser)

	server := http.Server{
		Handler: &serveMux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "  text/html")
	w.WriteHeader(200)
	html := fmt.Sprintf(`
	<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, cfg.fileserverHits.Load())
	bytes := fmt.Appendf(nil, "%s", html)
	w.Write(bytes)
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	_, err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	w.Header().Set("content-type", "  text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Reset called"))

}
