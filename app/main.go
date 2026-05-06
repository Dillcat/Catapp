package main

import (
	"net/http"
	"time"
	"log"
	"fmt"
	"sync/atomic"
	"encoding/json"
	"strings"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type chirpError struct {
	Error string `json:"error"`
}

type chirpValid struct {
	Cleaned_body string `json:"cleaned_body"`
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
	
	cleansedChirp := chirpValid{}

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
		cleansedChirp = wordCleanse(w, params)
		respondJson(w, 200, cleansedChirp)
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

func respondJson(w http.ResponseWriter, code int, cleansedChirp chirpValid){
	w.Header().Set("Content-Type", "application/json")


	data, err := json.Marshal(cleansedChirp)

	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong")
	}

	w.WriteHeader(200)
	w.Write(data)
}

func wordCleanse(w http.ResponseWriter, params parameters) (cleansedChirp chirpValid) {
	w.Header().Set("Content-Type", "application/json")

	body := params.Body

	
	curses := map[string]bool{
		"kerfuffle": true, 
		"sharbert": true,
		"fornax": true,
	}

	words := strings.Split(body, " ")

	for i, word := range words {
		word = strings.ToLower(word)
		if curses[word] {
			words[i] = "****"
			
		}
	}

	joinedwords := strings.Join(words, " ")

	cleansedChirp.Cleaned_body = joinedwords
	return cleansedChirp
}