package main

import "net/http"

func (cfg *apiConfig) resetUsers(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteUsers(r.Context())
}