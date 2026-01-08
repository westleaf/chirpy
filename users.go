package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/westleaf/chirpy/internal/auth"
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
		Password string `json:"password"`
		Email    string `json:"email"`
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

	if params.Password == "" {
		respondWithError(w, 403, "password is required")
		return
	}

	if params.Email == "" {
		respondWithError(w, 400, "email is required")
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "internal server error")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed,
	})
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

func (cfg *apiConfig) updateUserEmailPassword(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "no token found")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "invalid token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("error closing request body: %s", err)
		}
	}()

	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "incorrect request body")
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, 401, "missing paramater in body")
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, "could not hash password")
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashed,
	})
	if err != nil {
		respondWithError(w, 500, "error changing password or email")
		return
	}

	respondWithJSON(w, 200, UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
	})
}
