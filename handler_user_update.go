package main

import (
	"encoding/json"
	"net/http"

	"github.com/dmandevv/chirpy/internal/auth"
	"github.com/dmandevv/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
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

	if params.Email == "" {
		respWithError(w, http.StatusBadRequest, "Email cannot be empty", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	user, err := cfg.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	respWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})

}
