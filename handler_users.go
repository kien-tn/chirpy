package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func handlerUsers(apiCfg *apiConfig, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}
	// create a new user with apiCfg.db.CreateUser
	u, err := apiCfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         u.ID,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
		"email":      u.Email,
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
