package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kien-tn/chirpy/internal/auth"
	"github.com/kien-tn/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = 3600
	}
	// get the user with apiCfg.db.GetUserByEmail
	u, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting user", err)
		return
	}
	// check if the password is correct
	if err := auth.CheckPasswordHash(u.HashedPassword, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	token, err := auth.MakeJWT(u.ID, cfg.secretKey, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    u.ID,
		ExpiredAt: time.Now().Add(24 * 60 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		Email:        u.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authorization header required", err)
		return
	}
	rt, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}
	if rt.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", nil)
		return
	}
	if rt.ExpiredAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}
	token, err := auth.MakeJWT(rt.UserID, cfg.secretKey, 3600*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authorization header required", err)
		return
	}
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token", err)
		return
	}
	// respond with a 204 No Content
	w.WriteHeader(http.StatusNoContent)

}
