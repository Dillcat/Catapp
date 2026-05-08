package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong, getChirp")
		return
	}

	var chirps []Chirp

	for _, dbc := range dbChirps {
		
		chirps = append(chirps, Chirp{
			ID: dbc.ID,
			CreatedAt: dbc.CreatedAt,
			UpdatedAt: dbc.UpdatedAt,
			Body: dbc.Body,
			UserID: dbc.UserID,
		})
	}
	respondJson(w, 200, chirps)
}