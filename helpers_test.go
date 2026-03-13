package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/borisfritz/chirpy/internal/database"
	"github.com/google/uuid"
)

//NOTE: mockDB struct that matches internal/database package
type mockDB struct {
    createChirpFn   func(ctx context.Context, params database.CreateChirpParams) (database.Chirp, error)
    getAllChirpsFn  func(ctx context.Context) ([]database.Chirp, error)
    getChirpByIDFn  func(ctx context.Context, id uuid.UUID) (database.Chirp, error)
    createUserFn    func(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	getUserByEmailFn func(ctx context.Context, email string) (database.User, error)
	resetUsersFn	func(ctx context.Context) error
}
func (m *mockDB) CreateChirp(ctx context.Context, params database.CreateChirpParams) (database.Chirp, error) {
    return m.createChirpFn(ctx, params)
}
func (m *mockDB) GetAllChirps(ctx context.Context) ([]database.Chirp, error) {
    return m.getAllChirpsFn(ctx)
}
func (m *mockDB) GetChirpByID(ctx context.Context, id uuid.UUID) (database.Chirp, error) {
    return m.getChirpByIDFn(ctx, id)
}
func (m *mockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
    return m.createUserFn(ctx, arg)
}
func (m *mockDB) ResetUsers(ctx context.Context)  error {
	return m.resetUsersFn(ctx)
}

//NOTE: helper funcs for testing
// create mock database
func newMockDB() *mockDB {
    return &mockDB{
createChirpFn: func(ctx context.Context, params database.CreateChirpParams) (database.Chirp, error) {
            return mockChirp(), nil
        },
        getAllChirpsFn: func(ctx context.Context) ([]database.Chirp, error) {
            return []database.Chirp{mockChirp(), mockChirp()}, nil
        },
        getChirpByIDFn: func(ctx context.Context, id uuid.UUID) (database.Chirp, error) {
            return mockChirp(), nil
        },
        createUserFn: func(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
            return database.User{}, nil
        },
		getUserByEmailFn: func(ctx context.Context, email string) (database.User, error) {
			return mockUser(), nil
		},
        resetUsersFn: func(ctx context.Context) error {
            return nil
        },
    }
}

// Create mock chirp
func mockChirp() database.Chirp {
    return database.Chirp{
        ID:        uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Body:      "this is a valid chirp",
        UserID:    uuid.New(),
    }
}

// Create Mock User
func mockUser() database.User {
	return database.User{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email: "example@email.com",
		HashedPassword: "example",
	}
}

// create mock chirp request body
func chirpRequestBody(body string) string {
    return `{"body":"` + body + `","user_id":"` + uuid.New().String() + `"}`
}

// create mock requests
func newChirpRequest(t *testing.T, method, path, body string) (*httptest.ResponseRecorder, *http.Request) {
    t.Helper()
    req := httptest.NewRequest(method, path, strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    return w, req
}
