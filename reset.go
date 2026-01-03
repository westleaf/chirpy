package main

import (
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		respondWithError(w, 403, "incorrect platform, only allowed on dev")
		return
	}

	cfg.fileserverHits.Store(0)
	log.Printf("reset metrics")

	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, 400, "could not reset users")
		return
	}
	log.Printf("cleared users from database")

	respondWithJSON(w, 200, map[string]string{
		"body": "OK",
	})
}
