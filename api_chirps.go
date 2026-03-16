package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/borisfritz/chirpy/internal/auth"
	"github.com/borisfritz/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpResponse struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	Body 		string 	  	`json:"body"`
	UserID 		uuid.UUID 	`json:"user_id"`
}

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	
	// validate headers
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	//read body into params
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	//validate chirp (jwt and lenght)
	user, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to validate token")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	//filter chirp
	filtered := chirpFilter(params.Body)

	// create responce and respond
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: filtered,
		UserID: user,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create chirp")
		return
	}
	respondWithJSON(w, http.StatusCreated, chirpResponse{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to get chirps")
		return
	}
	response := []chirpResponse{}
	for _, chirp := range chirps {
		response = append(response, chirpResponse{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.CreatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}
	respondWithJSON(w, http.StatusOK, chirpResponse{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.CreatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func chirpFilter(body string) string {
	banned := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	for i, word := range words {
		if slices.Contains(banned, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
