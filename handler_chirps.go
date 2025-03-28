package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/kien-tn/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChip(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
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
	// Create a new chirp with apiCfg.db.CreateChirp
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		log.Printf("Error parsing UserID: %s", err)
		w.WriteHeader(400)
		w.Write([]byte(`{"error": "Invalid UserID format"}`))
		return
	}

	c, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})

	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         c.ID,
		"created_at": c.CreatedAt,
		"updated_at": c.UpdatedAt,
		"body":       c.Body,
		"user_id":    c.UserID,
	})
}
