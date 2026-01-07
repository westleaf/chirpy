package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/westleaf/chirpy/internal/auth"
	"github.com/westleaf/chirpy/internal/database"
)

type userLoginResponse struct {
	User         UserResponse
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshTokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
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

	expireTime := time.Hour

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expireTime)
	if err != nil {
		respondWithError(w, 401, "invalid token")
		return
	}

	rfToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 401, "could not create token")
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     rfToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, 401, "could not create token")
		return
	}

	respondWithJSON(w, 200, userLoginResponse{
		User: UserResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "invalid token")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "invalid token")
		return
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 401, "could not make token")
		return
	}

	respondWithJSON(w, 200, refreshTokenResponse{
		Token: jwt,
	})
}
