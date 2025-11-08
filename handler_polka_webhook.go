package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"database/sql"
	"errors"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	var params parameters

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err := cfg.db.UpdateUserChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "DB error", err)
		}
    return
	}

	w.WriteHeader(http.StatusNoContent)
}