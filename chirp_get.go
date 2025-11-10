package main

import (
	"net/http"
	"sort"
	"github.com/google/uuid"
	"github.com/nhatquang342/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")

	var (
		dbChirps []database.Chirp
		err error
	)

	if authorIDStr == "" {
		dbChirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	} else {
		authorID, convErr := uuid.Parse(authorIDStr)
		if convErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsOneAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	sortOption := r.URL.Query().Get("sort")
	sort.Slice(chirps, func(i, j int) bool {
		if sortOption == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})
	
	/*
	if sortOption == "asc" || sortOption == "" {
		sort.Slice(chirps, func(i, j int) bool {return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)})
	} else if sortOption == "desc" {
		sort.Slice(chirps, func(i, j int) bool {return chirps[i].CreatedAt.After(chirps[j].CreatedAt)})
	}
	*/

	respondWithJSON(w, http.StatusOK, chirps)
}