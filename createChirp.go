package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/nhatquang342/chirpy/internal/database"
)

type Chirp struct {
	ID 		  uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body 	  string 	`json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type createChirpParams struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	var newChirp createChirpParams
	if err := json.NewDecoder(r.Body).Decode(&newChirp); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if newChirp.Body == "" {
		http.Error(w, "Chirp cannot be empty", http.StatusBadRequest)
		return
	}
	const maxChirpLength = 140
	if len(newChirp.Body) > maxChirpLength {
		http.Error(w, "Chirp is too long", http.StatusBadRequest)
		return
	}
	_, err := cfg.db.GetUserByID(r.Context(), newChirp.UserID)
	if err != nil {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}
	cleanedMsg := msgCleaner(newChirp.Body)

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedMsg,
		UserID: newChirp.UserID,
	})
	if err != nil {
		http.Error(w, "System failed to create chirp", http.StatusInternalServerError)
		return
	}
	chirp := Chirp{
		ID: 	   dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: 	   dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	
	resp, err := json.Marshal(chirp)
	if err != nil {
		http.Error(w, "Failed to generate chirp", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func msgCleaner(rawText string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(rawText, " ")
	for i, w := range words {
		lower := strings.ToLower(w)
		for _, bad := range badWords {
			if lower == bad {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}