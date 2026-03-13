package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/borisfritz/chirpy/internal/auth"
	"github.com/borisfritz/chirpy/internal/database"
	"github.com/google/uuid"
)

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerPostUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	 string	`json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to hash password")
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, userResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password	string	`json:"password"`
		Email 		string 	`json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "user not found")
		return
	}
	ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to check password")
		return
	}
	if !ok {
		respondWithError(w, http.StatusForbidden, "incorrect password")
		return
	}
	respondWithJSON(w, http.StatusOK, userResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}
