package main

import (
	"net/http"
	"time"

	"github.com/dmandevv/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "No refresh token provided", err)
		return
	}
	dbUser, err := cfg.queries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	jwtToken, err := auth.MakeJWT(dbUser.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respWithError(w, http.StatusUnauthorized, "Failed to create JWT token", err)
		return
	}

	respWithJSON(w, http.StatusOK, response{
		Token: jwtToken,
	})
}
