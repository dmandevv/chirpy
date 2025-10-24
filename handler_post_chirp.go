package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/dmandevv/chirpy/internal/auth"

	"github.com/dmandevv/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

var profanityRegex = regexp.MustCompile(`(?i)\b(kerfuffle|sharbert|fornax)\b`)

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Body string `json:"body"`
	}

	//Validate user using JWT token

	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respWithError(w, http.StatusUnauthorized, "Missing authorization token", err)
		return
	}

	userID, err := auth.ValidateJWT(jwtToken, cfg.JWTSecret)
	if err != nil {
		respWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// Decode request body

	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err = decoder.Decode(&params)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "Invalid request payload", err)
		return
	}

	// Validate chirp

	cleanedChirp, err := validateChirp(params.Body)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "Invalid chirp content: "+err.Error(), err)
		return
	}

	// Create chirp in database

	chirp, err := cfg.queries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: userID,
	})
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "Failed to create chirp", err)
		return
	}

	respWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func validateChirp(body string) (string, error) {

	if len(body) == 0 {
		return "", errors.New("Chirp is empty")
	}

	if len(body) > 140 {
		return "", errors.New("Chirp is too long")
	}

	return profanityFilter(body), nil
}

func profanityFilter(body string) string {
	return profanityRegex.ReplaceAllString(body, "****")
}
