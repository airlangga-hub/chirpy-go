package main

import (
	"net/http"
	"github.com/google/uuid"
	"sort"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorID := uuid.Nil
	var err error
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}

	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
	}

	var chirps []Chirp

	for _, dbChirp := range dbChirps{
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}

		chirps = append(chirps, Chirp{
			dbChirp.ID,
			dbChirp.CreatedAt,
			dbChirp.UpdatedAt,
			dbChirp.Body,
			dbChirp.UserID,
		})
	}

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "desc" {
		sort.Slice(
			chirps,
			func(i, j int) bool {return chirps[i].CreatedAt.After(chirps[j].CreatedAt)},
		)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirpDb, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
	return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		chirpDb.ID,
		chirpDb.CreatedAt,
		chirpDb.UpdatedAt,
		chirpDb.Body,
		chirpDb.UserID,
	})
}