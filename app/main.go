package main

import (
	"net/http"
	"time"
	"log"
	"fmt"
)


func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr: ":8080", 
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		
		// sets Header map to key Content-Type, value text/plain...
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// writes status code 200
		w.WriteHeader(http.StatusOK)

		// accepts []bytes, returns int, error
		_, err := w.Write([]byte("OK"))
		if err != nil {
			fmt.Printf("%v", err)
		}


	})
	log.Fatal(server.ListenAndServe())

}

