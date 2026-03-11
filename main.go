package main

import (
	"os"
	"log"
	"database/sql"
	"net/http"
	"sync/atomic"

	"github.com/borisfritz/chirpy/internal/database"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits 	atomic.Int32
	db 				*database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	dbConn , err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		return
	}
	dbQueries := database.New(dbConn)

	cfg := &apiConfig{db: dbQueries}
	mux := http.NewServeMux()
	
	//app handlers
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	
	//api handlers
	//GET Requests only
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	//POST Requests only
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	// Create and start server
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
