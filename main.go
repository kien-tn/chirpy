package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment the fileserverHits counter
		cfg.fileserverHits.Add(1)

		// Print the current hit count to stdout
		fmt.Fprintln(os.Stdout, "Hitting:", cfg.fileserverHits.Load())

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	apiCfg := &apiConfig{}
	fmt.Fprintln(os.Stdout, "Hitting:", apiCfg.fileserverHits.Load())
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", middlewareLog(apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("."))))))
	mux.Handle("/healthz", middlewareLog(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ContentType
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})))
	mux.Handle("/metrics", middlewareLog(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ContentType
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits: " + strconv.Itoa(int(apiCfg.fileserverHits.Load()))))
	})))
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()

}
