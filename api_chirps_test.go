package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newChirpRequest(t *testing.T, body string) (*httptest.ResponseRecorder, *http.Request) {
    t.Helper()
    req := httptest.NewRequest("POST", "/api/validate_chirp", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    return w, req
}

func TestValidateChirpValid(t *testing.T) {
	w, req := newChirpRequest(t, `{"body":"this is a valid chirp"}`)
    handlerValidateChirp(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }
}

func TestValidateChirpTooLong(t *testing.T) {
	w, req := newChirpRequest(t, `"body": "`+strings.Repeat("a", 141)+`"`)
    handlerValidateChirp(w, req)
    if w.Code != http.StatusBadRequest {
        t.Errorf("expected 400, got %d", w.Code)
    }
}

func TestValidateChirpInvalidJSON(t *testing.T) {
	w, req := newChirpRequest(t, `not json`)
    handlerValidateChirp(w, req)
    if w.Code != http.StatusBadRequest {
        t.Errorf("expected 400, got %d", w.Code)
    }
}

func TestChirpFilter(t *testing.T) {
    input := "this is a kerfuffle"
    expected := "this is a ****"
    result := chirpFilter(input)
    if result != expected {
        t.Errorf("expected %s, got %s", expected, result)
    }
}
