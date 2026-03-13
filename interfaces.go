package main

import (
	"context"
	
	"github.com/google/uuid"
	"github.com/borisfritz/chirpy/internal/database"
)

type DB interface {
	CreateChirp(ctx context.Context, params database.CreateChirpParams) (database.Chirp, error)
	GetAllChirps(ctx context.Context) ([]database.Chirp, error)
	GetChirpByID(ctx context.Context, id uuid.UUID) (database.Chirp, error)
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUserByEmail(ctx context.Context, email string) (database.User, error)
	ResetUsers(ctx context.Context) error
}
