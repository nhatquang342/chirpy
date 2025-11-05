package main

import "net/http"

func (cfg *apiConfig) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		http.Error(w, "Deleting outside dev environment is not allowed", http.StatusForbidden)
		return
	}
	
	cfg.fileserverHits.Store(0)
	
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to reset user base", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users have been deleted"))
}