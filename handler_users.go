package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kien-tn/chirpy/internal/auth"
	"github.com/kien-tn/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func handlerUsers(apiCfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}
	if params.Password == "" || params.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Email and password are required"}`))
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}
	// create a new user with apiCfg.db.CreateUser
	u, err := apiCfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            u.ID,
		"created_at":    u.CreatedAt,
		"updated_at":    u.UpdatedAt,
		"email":         u.Email,
		"is_chirpy_red": u.IsChirpyRed,
	})
}

func handlerUsersReset(apiCfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	env := os.Getenv("PLATFORM")
	if env != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "Resetting the database is only allowed in dev environment"}`))
		return
	}
	err := apiCfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error resetting database: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters", err)
		return
	}
	if params.Password == "" || params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}
	// Update user with apiCfg.db.CreateUser
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	u, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerUpdateUserRed(w http.ResponseWriter, r *http.Request) {
	type inner struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  inner  `json:"data"`
	}
	defer r.Body.Close()
	//Check API key
	polkaKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}
	if polkaKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	} else {
		if params.Data.UserID == uuid.Nil {
			respondWithError(w, http.StatusNotFound, "User ID is required", nil)
			return
		}
		u, err := cfg.db.UpdateUserChirpyRed(r.Context(), params.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
			return
		}
		respondWithJSON(w, http.StatusNoContent, map[string]interface{}{
			"id":            u.ID,
			"is_chirpy_red": u.IsChirpyRed,
		})
	}

}
