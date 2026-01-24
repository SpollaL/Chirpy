package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/SpollaL/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	platform       string
	secret_key     string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret_key := os.Getenv("SECRET_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not open connection to database")
		return
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		queries:        dbQueries,
		platform:       platform,
		secret_key:     secret_key,
	}
	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))),
	)
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandleChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandleGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandleGetChirp)
	mux.HandleFunc("POST /api/users", apiCfg.HandleUserCreation)
	mux.HandleFunc("PUT /api/users", apiCfg.HandleUserUpdate)
	mux.HandleFunc("POST /api/login", apiCfg.HandleUserLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandleTokenRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandleTokenRevoke)
	server := http.Server{Handler: mux, Addr: ":8080"}
	server.ListenAndServe()
}
