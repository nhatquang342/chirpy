package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// next http.Handler is the actual handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // a new handler that includes extra behavior
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r) // forwards the request to the real handler
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`, cfg.fileserverHits.Load()))) // Load() is to read the atomic counter safely
}