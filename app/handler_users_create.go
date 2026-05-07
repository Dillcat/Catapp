package main

import (
	"encoding/json"
	"net/http"
	"time"
	"log"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (cfg *apiConfig)createUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	user := User{}

	err := decoder.Decode(&user)
	
	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong")
	}
	//CreateUser is an sqlc function that populates the user struct
	u, err := cfg.db.CreateUser(r.Context(), user.Email)

	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong, createUser")
		return
	}

	userdata := User{
		ID: u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email: u.Email,
	}

	respondJsonUser(w, userdata)
}

func respondJsonUser(w http.ResponseWriter, userdata User) {
	w.Header().Set("Content-Type", "application/json")


	data, err := json.Marshal(userdata)

	if err != nil {
		log.Printf("Error", err)
		respondError(w, 500, "Something went wrong, Json Encoding")
		return
	}

	w.WriteHeader(201)
	w.Write(data)
}
