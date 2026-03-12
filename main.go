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
	platform		string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	dbConn , err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
    	log.Fatal("Unable to connect to database: ", err)
	}

	dbQueries := database.New(dbConn)

	cfg := &apiConfig{
		db: dbQueries, 
		platform: os.Getenv("PLATFORM"),
	}
	mux := http.NewServeMux()
	
	//app handlers
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	
	//api handlers
	//GET Requests
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirp)
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	//POST Requests
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/users", cfg.handlerPostUsers)
	mux.HandleFunc("POST /api/chirps", cfg.handlerPostChirp)

	// Create and start server
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
