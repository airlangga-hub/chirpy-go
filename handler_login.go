package main

import (
	"net/http"
	"encoding/json"
	"github.com/airlangga-hub/chirpy-go/internal/auth"
	"time"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	var params parameters

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid user email", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !match || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthrozied user", err)
		return
	}

	expiresIn := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expiresIn = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	access_token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			user.ID,
			user.CreatedAt,
			user.UpdatedAt,
			user.Email,
		},
		Token: access_token,
	})
}