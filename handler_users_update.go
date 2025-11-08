package main

import (
	"net/http"
	"github.com/airlangga-hub/chirpy-go/internal/auth"
	"github.com/airlangga-hub/chirpy-go/internal/database"
	"encoding/json"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find access token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate access token", err)
		return
	}

	var params parameters

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	userUpdated, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
		ID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: userUpdated.ID,
		CreatedAt: userUpdated.CreatedAt,
		UpdatedAt: userUpdated.UpdatedAt,
		Email: userUpdated.Email,
		IsChirpyRed: userUpdated.IsChirpyRed,
	})
}