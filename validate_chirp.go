package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validate_chirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Failed to parse body. Invalid JSON")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	words := strings.Split(params.Body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if loweredWord == "kerfuffle" || loweredWord == "sharbert" || loweredWord == "fornax" {
			words[i] = "****"
		}
	}
	sanitazedChirp := strings.Join(words, " ")
	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}
	response := responseBody{
		CleanedBody: sanitazedChirp,
	}
	respondWithJSON(w, 200, response)
}
