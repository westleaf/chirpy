package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/westleaf/chirpy/internal/auth"
	"github.com/westleaf/chirpy/internal/database"
)

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "unauthorized")
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		return
	}

	if params.Body == "" {
		respondWithError(w, 400, "body can not be empty")
		return
	}

	if userId == uuid.Nil {
		respondWithError(w, 400, "user_id can not be empty")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleaned := CensorChirp(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, 500, "could not create chirp")
		return
	}

	respondWithJSON(w, 201, ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "could not get chirps")
		return
	}

	responses := make([]ChirpResponse, 0, len(chirps))

	for _, chirp := range chirps {
		responses = append(responses, ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}

	respondWithJSON(w, 200, responses)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpUuid, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, 500, "could not parse chirp id")
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpUuid)
	if err != nil {
		respondWithError(w, 404, "could not find chirp")
	}

	respondWithJSON(w, 200, ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func CensorChirp(s string) string {
	banned := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Split(s, " ")

	for i, w := range words {
		if banned[strings.ToLower(w)] {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
