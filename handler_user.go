package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dmandevv/chirpy/internal/auth"
	"github.com/dmandevv/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Password    string    `json:"password,omitempty"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if req.Email == "" {
		respWithError(w, http.StatusBadRequest, "Email is required", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	respWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
