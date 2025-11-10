package main

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"github.com/nhatquang342/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpdateMembership(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	APIKey, err := auth.GetAPIKey(r.Header)
	if err != nil || APIKey != cfg.polka_key {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpdateMembership(r.Context(), params.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}