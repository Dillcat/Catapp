package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"log"
)

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {
	
	cleansedChirp := chirpValid{}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondError(w, 400, "Chirp is too long")
		return
	} else {
		cleansedChirp = wordCleanse(w, params)
		respondJson(w, 200, cleansedChirp)
	}


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