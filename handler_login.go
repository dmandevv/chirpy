package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dmandevv/chirpy/internal/auth"
	"github.com/dmandevv/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if !match || err != nil {
		respWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	//Create JWT token
	jwtToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "Failed to create JWT token", err)
		return
	}

	//Create refresh token
	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "Failed to create refresh token", err)
		return
	}

	respWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        jwtToken,
		RefreshToken: refreshToken,
	})
}
