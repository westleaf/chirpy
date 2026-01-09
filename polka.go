package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/westleaf/chirpy/internal/auth"
)

func (cfg *apiConfig) upgradeUserToRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	_, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "invalid key")
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
		respondWithError(w, 400, "invalid request body")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUserToRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
