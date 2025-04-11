package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kien-tn/chirpy/internal/auth"
	"github.com/kien-tn/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerCreateChip(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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
	token, _ := auth.GetBearerToken(r.Header)
	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token, missing UserID", err)
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

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	output := []Chirp{}
	var chirps []database.Chirp
	var err error
	s := r.URL.Query().Get("author_id")
	log.Println("author_id found in request path: ", s)
	if s != "" {
		authorID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID format", err)
			return
		}
		chirps, err = cfg.db.GetChirpsByUserID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error getting chirps", err)
			return
		}
	} else {
		chirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error getting chirps", err)
			return
		}
	}
	for _, c := range chirps {
		output = append(output, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, output)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirp_id")
	log.Println("chirp ID found in request path: ", id)
	chirpID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
		return
	}
	c, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error getting chirp", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirp_id")
	log.Println("chirp ID found in request path: ", id)
	chirpID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token, missing UserID", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error getting chirp", err)
		return
	}
	// Check if the chirp belongs to the user
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You do not have permission to delete this chirp", nil)
		return
	}
	err = cfg.db.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{
		ID:     chirpID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
