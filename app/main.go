package main

import _ "github.com/lib/pq"

import (
	"net/http"
	"time"
	"log"
	"sync/atomic"
	"encoding/json"
	"os"
	"github.com/Dillcat/Catapp/app/internal/database"
	"github.com/joho/godotenv"
	"database/sql"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

type chirpError struct {
	Error string `json:"error"`
}

type chirpValid struct {
	Cleaned_body string `json:"cleaned_body"`
}

//type parameters struct {
//	Body string `json:"body"`
//}


func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{}

	apiCfg.db = dbQueries

	mux := http.NewServeMux()
	
	//mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)

	mux.HandleFunc("POST /admin/reset", apiCfg.resetUsers)

	mux.HandleFunc("POST /api/chirps", apiCfg.validateChirp)

	mux.HandleFunc("POST /api/users", apiCfg.createUser)

	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)

	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)

	server := &http.Server{
		Addr: ":8080", 
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())

}

func respondError(w http.ResponseWriter, code int, msg string){
	w.Header().Set("Content-Type", "application/json")

	chirpErr := chirpError{}

	chirpErr.Error = msg

	data, err := json.Marshal(chirpErr)
	
	if err != nil {
		log.Printf("Error", err)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func respondJson(w http.ResponseWriter, code int, payload interface{}){
	w.Header().Set("Content-Type", "application/json")


	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong")
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

