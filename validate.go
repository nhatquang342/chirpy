package main

import (
	"encoding/json"
	"log"
	"net/http"
	//"strings" as msgCleaner is now commented out
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned string `json:"cleaned_body"`
	}
	
	decoder := json.NewDecoder(r.Body)
	newChirp := chirp{}
	err := decoder.Decode(&newChirp)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Error decoding request body", err)
		return
	}
	
	const maxChirpLength = 140
	if len(newChirp.Body) > maxChirpLength {
		respondWithErr(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedMsg := msgCleaner(newChirp.Body)
	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned: cleanedMsg,
	})
}
/*
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
*/
func respondWithErr(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code >= 500 {
		log.Printf("Error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(resp)
}