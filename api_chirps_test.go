package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/borisfritz/chirpy/internal/database"
	"github.com/google/uuid"
)

//NOTE: Unit Tests
//NOTE: handlerPostChirp Tests
func TestPostChirpValid(t *testing.T) {
    // 1. create a default mock db
	cfg := &apiConfig{db: newMockDB()}
    // 2. create a fake request with a valid JSON body
    w, req := newChirpRequest(t, "POST", "/api/chirps", chirpRequestBody("valid chirp"))
    // 3. call the handler directly, passing the fake request and response recorder
    cfg.handlerPostChirp(w, req)
    // 4. assert the response code is what we expect
    if w.Code != http.StatusCreated {
        t.Errorf("expecting 201, got %d", w.Code)
    }
}

func TestPostChirpTooLong(t *testing.T) {
	cfg := &apiConfig{db: newMockDB()}
	w, req := newChirpRequest(t, "POST", "/api/chirps", chirpRequestBody(strings.Repeat("a",141)))
	cfg.handlerPostChirp(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expecting 400, got %v", w.Code)
	}
}

func TestPostChirpInvalidJSON(t *testing.T) {
	cfg := &apiConfig{db: newMockDB()}
	w, req := newChirpRequest(t, "POST", "/api/chirps", `not JSON`)
	cfg.handlerPostChirp(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expecting 400, got %v", w.Code)
	}
}

func TestPostChirpDBError(t *testing.T) {
	db := newMockDB()
	db.createChirpFn = func(ctx context.Context, params database.CreateChirpParams) (database.Chirp, error) {
		return database.Chirp{}, errors.New("db error")
	}
	cfg := &apiConfig{db: db}
	w, req := newChirpRequest(t, "POST", "/api/chirps", chirpRequestBody("valid chirp"))
	cfg.handlerPostChirp(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expecting 500, got %v", w.Code)
	}
}

//NOTE: handlerGetChirps tests
func TestGetChirpsValid(t *testing.T) {
	cfg := &apiConfig{db: newMockDB()}
	w, req := newChirpRequest(t, "GET", "/api/chirps", "")
	cfg.handlerGetChirps(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expecting 200, got %v", w.Code)
	}
}

func TestGetChirpsEmpty(t *testing.T) {
	db := newMockDB()
	db.getAllChirpsFn = func(ctx context.Context) ([]database.Chirp, error) {
		return []database.Chirp{}, nil
	}
	cfg := &apiConfig{db: db}
	w, req := newChirpRequest(t, "GET", "/api/chirps", "")
	cfg.handlerGetChirps(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expecting 200, got %v", w.Code)
	}
}

func TestGetChirpsDBError(t *testing.T) {
	db := newMockDB()
	db.getAllChirpsFn = func(ctx context.Context) ([]database.Chirp, error) {
		return []database.Chirp{}, errors.New("db error")
	}
	cfg := &apiConfig{db: db}
	w, req := newChirpRequest(t, "GET", "/api/chirps", "")
	cfg.handlerGetChirps(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expecting 500, got %v", w.Code)
	}
}

//NOTE: handlerGetChirp tests
func TestGetChirpValid(t *testing.T) {
	chirp := mockChirp()
	db := newMockDB()
	db.getChirpByIDFn = func(ctx context.Context, id uuid.UUID) (database.Chirp, error) {
		return chirp, nil
	}
	cfg := &apiConfig{db: db}
	w, req := newChirpRequest(t, "GET", "/api/chirps/"+chirp.ID.String(), "")
	req.SetPathValue("chirpID", chirp.ID.String())
	cfg.handlerGetChirp(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expecting 200, got %v", w.Code)
	}
}

func TestGetChirpInvalid(t *testing.T) {
	cfg := &apiConfig{db: newMockDB()}
	w, req := newChirpRequest(t, "GET", "/api/chirps/not-a-uuid", "")
	req.SetPathValue("chirpID", "not-a-uuid")
	cfg.handlerGetChirp(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expecting 400, got %v", w.Code)
	}
}

func TestGetChirpNotFound(t *testing.T) {
	db := newMockDB()
	db.getChirpByIDFn = func(ctx context.Context, id uuid.UUID) (database.Chirp, error) {
		return database.Chirp{}, errors.New("not found")
	}
	cfg := &apiConfig{db: db}
	w, req := newChirpRequest(t, "GET", "/api/chirps/"+uuid.New().String(), "")
	req.SetPathValue("chirpID", uuid.New().String())
	cfg.handlerGetChirp(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expecting 404, got %v", w.Code)
	}
} 

// chirpFilter tests
func TestChirpFilter(t *testing.T) {
    input := "this is a kerfuffle"
    expected := "this is a ****"
    result := chirpFilter(input)
    if result != expected {
        t.Errorf("expected %s, got %s", expected, result)
    }
}
