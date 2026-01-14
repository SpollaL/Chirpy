package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/SpollaL/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not open connection to database")
		return
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}, queries: dbQueries, platform: platform}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)
	mux.HandleFunc("POST /api/validate_chirp", HandleChirpValidation)
	mux.HandleFunc("POST /api/users", apiCfg.HandleUserCreation)
	server := http.Server{Handler: mux, Addr: ":8080"}
	server.ListenAndServe()
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits atomic.Int32
	queries *database.Queries
	platform string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) HandleMetrics(w http.ResponseWriter, _ *http.Request) {
	content := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
}

func (cfg *apiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
	}
	err := cfg.queries.DeleteUsers(r.Context())	
	if err != nil {
		log.Fatalf("Could not delete all users: %v", err)
	}
	cfg.fileserverHits.Store(0)
}

func HandleChirpValidation(w http.ResponseWriter, req *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)

	w.Header().Set("Content-Type", "text/html")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "Unable to decode json body"}`))
		return
	}

	if len(chirp.Body) > 140 {
		w.WriteHeader(400)
		w.Write([]byte(`{"error": "Chirp is too long. Should be less than 140 characters"}`))
		return
	}

	chirp.Body = ReplaceProfane(chirp.Body)
	cleaned_body := fmt.Sprintf(`{"cleaned_body": "%s"}`, chirp.Body)

	w.WriteHeader(200)
	w.Write([]byte(cleaned_body))
}

func ReplaceProfane(s string) string {
	profane_words := []string{"kerfuffle", "sharbert", "fornax"}
	for word:= range strings.SplitSeq(s, " ") {
		for _, pword:= range profane_words {
			if strings.ToLower(word) == pword {
				splits := strings.Split(s, word)
				s = strings.Join(splits, "****")
			}
		}
	}
	return s 
}

func (cfg *apiConfig) HandleUserCreation(w http.ResponseWriter, r *http.Request) {
	type ReqStruct struct {
		Email string `json:"email"`
	}	
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
  }
	decoder := json.NewDecoder(r.Body)
	reqStruct := &ReqStruct{}
	err := decoder.Decode(reqStruct)
	if err != nil {
		log.Fatalf("Could not decode json request: %v", err)
	}
	dbUser, err := cfg.queries.CreateUser(r.Context(), reqStruct.Email)
	jsonUser := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	}
	res, err := json.Marshal(jsonUser)
	if err != nil {
		log.Fatalf("Could not marshal json response: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(res)
}
