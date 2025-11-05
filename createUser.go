package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID 			uuid.UUID `json:"id"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	Email 		string 	  `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	type createUserParams struct {
		Email string `json:"email"`
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

	dbUser, err := cfg.db.CreateUser(r.Context(), newUser.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := User{
		ID: 	   dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: 	   dbUser.Email,
	}

	resp, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		http.Error(w, "Failed to create/encode user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	w.Write(resp)
}