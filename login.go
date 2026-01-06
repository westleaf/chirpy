package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/westleaf/chirpy/internal/auth"
)

type userLoginResponse struct {
	User  UserResponse
	Token string `json:"token"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
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

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "unauthorized")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, 401, "unauthorized")
		return
	}

	expireTime := time.Hour // default

	maxSeconds := int(time.Hour.Seconds())

	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < maxSeconds {
		expireTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expireTime)
	if err != nil {
		respondWithError(w, 401, "invalid token")
		return
	}

	respondWithJSON(w, 200, userLoginResponse{
		User: UserResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}
