package main

import (
	"net/http"

	"github.com/westleaf/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "error reading token")
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "error getting token")
		return
	}

	respondWithJSON(w, 204, http.NoBody)
}
