package main

import (
	"encoding/json"
	"net/http"
	"github.com/nhatquang342/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "User is not found", http.StatusUnauthorized)
		return
	}
	
	pwMatch, err := auth.CheckPasswordHash(req.Password, dbUser.HashedPassword)
	if err != nil || !pwMatch {
		http.Error(w, "Password incorrect", http.StatusUnauthorized)
		return
	} else {
		respondWithJSON(w, http.StatusOK, User{
			ID: 	   dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email: 	   dbUser.Email,
		})
	}

	
}