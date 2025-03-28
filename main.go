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

	"github.com/joho/godotenv"
	"github.com/kien-tn/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
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

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body        string `json:"body"`
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		// Add a message to the response body "error": "Something went wrong"
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	if len(params.Body) > 140 {
		// If the body is too long, return a 400 Bad Request
		w.WriteHeader(400)
		// Add a message to the response body "error": "Chirp is too long"
		w.Write([]byte(`{"error": "Chirp is too long"}`))
		return
	}
	// check if the body contains any words such as "kerfuffle" "sharbert" "fornax"
	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax", "Kerfuffle", "Sharbert", "Fornax"}
	_, cleanedBody := maskForbiddenWord(params.Body, forbiddenWords)
	response := map[string]interface{}{
		"valid": true,
	}
	response["cleaned_body"] = cleanedBody
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func maskForbiddenWord(body string, forbiddenWords []string) (bool, string) {
	bodyLower := strings.ToLower(body)
	containsForbiddenWord := false
	for _, word := range forbiddenWords {
		if strings.Contains(bodyLower, strings.ToLower(word)) {
			body = strings.ReplaceAll(body, word, "****")
			containsForbiddenWord = true
		}
	}
	return containsForbiddenWord, body
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)
	apiCfg := &apiConfig{
		db: dbQueries,
	}
	defer db.Close()
	fmt.Fprintln(os.Stdout, "Hitting:", apiCfg.fileserverHits.Load())
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", middlewareLog(apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("."))))))
	mux.Handle("GET /api/healthz", middlewareLog(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ContentType
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})))
	mux.Handle("GET /admin/metrics", middlewareLog(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ContentType
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`
<html>
  <body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
  </body>
</html>`, apiCfg.fileserverHits.Load())))
	})))
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		// apiCfg.fileserverHits.Store(0)
		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("OK"))
		handlerUsersReset(apiCfg, w, r)
	})
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		handlerUsers(apiCfg, w, r)
	})
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChip)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()

}
