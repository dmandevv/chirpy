package main

import (
	"net/http"

	"github.com/dmandevv/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respWithError(w, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.JWTSecret)
	if err != nil {
		respWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	chirp, err := cfg.queries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if chirp.UserID != userID {
		respWithError(w, http.StatusForbidden, "You are not allowed to delete this chirp", nil)
		return
	}

	err = cfg.queries.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "Failed to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
