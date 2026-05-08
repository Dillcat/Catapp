package main

import (
	"log"
	"net/http"
	"fmt"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(r.PathValue("chirpID"))
	id, err := uuid.Parse(r.PathValue("chirpID"))
	
	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong, parse")
		return
	}
	
	dbc, err := cfg.db.GetChirp(r.Context(), id)

	if err != nil {
		log.Printf("Error", err)
		respondError(w, 404, "")
		return
	}

	chirp := Chirp{
		ID: dbc.ID,
		CreatedAt: dbc.CreatedAt,
		UpdatedAt: dbc.UpdatedAt,
		Body: dbc.Body,
		UserID: dbc.UserID,
	}
	respondJson(w, 200, chirp)
}