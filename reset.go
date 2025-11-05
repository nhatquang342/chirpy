package main

import "net/http"

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // bootdev solution doesn't have this line
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to zero"))
}