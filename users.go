package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/westleaf/chirpy/internal/database"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("error closing request body: %s", err)
		}
	}()

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "invalid request body")
		return
	}

	if params.Email == "" {
		respondWithError(w, 400, "email is required")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("%s\n", err)
		respondWithError(w, 500, "could not create user")
		return
	}

	respondWithJSON(w, 201, NewUserResponse(user))
}

func NewUserResponse(u database.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}
