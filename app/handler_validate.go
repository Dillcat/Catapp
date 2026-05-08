package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"log"
	"github.com/google/uuid"
	"time"
	"github.com/Dillcat/Catapp/app/internal/database"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}

	err := decoder.Decode(&chirp)

	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong")
		return
	}

	if len(chirp.Body) > 140 {
		respondError(w, 400, "Chirp is too long")
		return
	} else {



		chirp = wordCleanse(w, chirp)

		c, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			Body: chirp.Body,
			UserID: chirp.UserID,
		})

		if err != nil {
			log.Printf("Error", err)
			respondError(w, 500, "Something went wrong, createChirp")
			return
		}

		chirp = Chirp{
			ID: c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body: c.Body,
			UserID: c.UserID,
		}

		respondJson(w, 201, chirp)
	}



}

func wordCleanse(w http.ResponseWriter, chirp Chirp) (cleansedChirp Chirp) {
	w.Header().Set("Content-Type", "application/json")

	body := chirp.Body

	
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

	chirp.Body = joinedwords
	return chirp
}