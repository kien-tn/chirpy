package main

import (
	"encoding/json"
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if len(params.Body) > 140 {
		// If the body is too long, return a 400 Bad Request
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	// Create a new chirp with apiCfg.db.CreateChirp
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UserID format", err)
		return
	}

	c, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
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
