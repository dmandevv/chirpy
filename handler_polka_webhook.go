package main

import (
	"encoding/json"
	"net/http"

	"github.com/dmandevv/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}
	if key != cfg.PolkaKey {
		respWithError(w, http.StatusUnauthorized, "API keys don't match", nil)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "Invalid user id format", err)
		return
	}

	_, err = cfg.queries.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		respWithError(w, http.StatusNotFound, "Failed to upgrade user to Chirpy Red", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
