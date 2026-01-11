package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/healthz/", healthz)
	mux.HandleFunc("/metrics/", apiCfg.metrics)
	mux.HandleFunc("/reset/", apiCfg.reset)
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
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, _ *http.Request) {
	content := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(content))
}

func (cfg *apiConfig) reset(_ http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
}
