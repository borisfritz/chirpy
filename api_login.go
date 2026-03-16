package main

import (
	"time"
	"net/http"
	"encoding/json"

	"github.com/borisfritz/chirpy/internal/auth"
	"github.com/borisfritz/chirpy/internal/database"
)


func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password			string	`json:"password"`
		Email 				string 	`json:"email"`
	}

	//Read into params
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	//auth user
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user not found", err)
		return
	}
	ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to check password", err)
		return
	}
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "incorrect password", err)
		return
	}

	//create tokens
	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create session token", err)
		return
	}
	refreshKey := auth.MakeRefreshToken()
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshKey,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create refrsh token", err)
		return
	}
	
	//respond
	respondWithJSON(w, http.StatusOK, userResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: jwtToken,
		RefreshToken: refreshToken.Token,
	})
}
