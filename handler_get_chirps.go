package main

import (
	"net/http"
	"slices"

	"github.com/dmandevv/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	chirpID := r.PathValue("chirpID")

	uID, err := uuid.Parse(chirpID)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}

	dbChirp, err := cfg.queries.GetChirpByID(r.Context(), uID)
	if err != nil {
		respWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	}

	respWithJSON(w, http.StatusOK, chirp)

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	var err error
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}

	dbChirps := []database.Chirp{}

	if authorID != uuid.Nil {
		dbChirps, err = cfg.queries.GetAllChirpsOfUser(r.Context(), authorID)
		if err != nil {
			respWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
			return
		}
	} else {
		dbChirps, err = cfg.queries.GetAllChirps(r.Context())
		if err != nil {
			respWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
			return
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	sort := r.URL.Query().Get("sort")
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	slices.SortFunc(chirps, func(i, j Chirp) int {
		if sort == "asc" {
			return i.CreatedAt.Compare(j.CreatedAt)
		}
		return j.CreatedAt.Compare(i.CreatedAt)
	})

	respWithJSON(w, http.StatusOK, chirps)

}
