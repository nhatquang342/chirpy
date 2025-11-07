package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/nhatquang342/chirpy/internal/auth"
	"github.com/nhatquang342/chirpy/internal/database"
)

type User struct {
	ID 			uuid.UUID `json:"id"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	Email 		string 	  `json:"email"`
	Token		string	  `json:"token"`
	RefreshToken string   `json:"refresh_token`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type createUserParams struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	var newUser createUserParams // newUser := createUserParams{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if newUser.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: newUser.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: 	   dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: 	   dbUser.Email,
	})
}