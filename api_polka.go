package main

import (
	"os"
	"strings"
	"errors"
	"encoding/json"
	"net/http"

	"github.com/borisfritz/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhooks (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey := os.Getenv("POLKA_KEY")
	headerAPIKey, err := GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to get api key", err)
		return
	}
	if apiKey != headerAPIKey {
		respondWithError(w, http.StatusUnauthorized, "api key mismatch", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse uuid", err)
		return
	}
	_, err = cfg.db.UpdateChirpyRed(r.Context(), database.UpdateChirpyRedParams{
		ID: userID,
		IsChirpyRed: true,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetAPIKey(h http.Header) (string, error) {
	authHeader := h.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No Authorization Header")
	}
	key := strings.TrimPrefix(authHeader, "ApiKey")
	if key == authHeader {
		return "", errors.New("invalid authorization header format")
	}
	return strings.TrimSpace(key), nil
} 
