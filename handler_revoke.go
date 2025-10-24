package main

import (
	"net/http"

	"github.com/dmandevv/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	revokeToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "Couldn't validate token", err)
		return
	}

	err = cfg.queries.RevokeRefreshToken(r.Context(), revokeToken)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
