package main

import (
	"net/http"
	"time"
	"log"
	"fmt"
	"sync/atomic"
	"encoding/json"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type chirpError struct {
	Error string `json:"error"`
}

type chirpValid struct {
	Valid bool `json:"valid"`
}

type parameters struct {
	Body string `json:"body"`
}


func main() {

	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	
	//mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)

	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)

	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirp)

	server := &http.Server{
		Addr: ":8080", 
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	val := fmt.Sprintf(adminstring, cfg.fileserverHits.Load())
	w.Write([]byte(val))
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong")
	}

	if len(params.Body) > 140 {
		respondError(w, 400, "Chirp is too long")
		return
	} else {
		respondJson(w, 200)
	}


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

func respondJson(w http.ResponseWriter, code int){
	w.Header().Set("Content-Type", "application/json")

	validChirp := chirpValid{}

	validChirp.Valid = true

	data, err := json.Marshal(validChirp)

	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong")
	}

	w.WriteHeader(200)
	w.Write(data)
}