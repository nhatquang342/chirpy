package main

import (
	"encoding/json"
	"net/http"
	"time"
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
	}	
	
	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 60)
	_, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		ExpiresAt: expiresAt,
		UserID: dbUser.ID,
	})
	if err != nil {
		http.Error(w, "Fail to create refresh token", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: 	   dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: 	   dbUser.Email,
		Token:	   accessTokenoken,
		RefreshToken: refreshToken,
	})
	
}