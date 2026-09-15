package main

import (
	"encoding/json"
	"net/http"
)

func chirp(responseWriter http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(responseWriter, 400, "Failed to, parse body. Invalid JSON")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(responseWriter, 400, "Chirp is too long")
	}

	type returnValid struct {
		Valid bool `json:"valid"`
	}
	validResponse := returnValid{
		Valid: true,
	}
	respondWithJSON(responseWriter, 200, validResponse)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}
	error := returnError{
		Error: msg,
	}
	dat, marshalErr := json.Marshal(error)
	if marshalErr != nil {
		w.WriteHeader(500)
		bytes := []byte("Failed to parse error.. internal server error")
		w.Write(bytes)
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(dat))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		bytes := []byte("Failed to parse payload... internal server error")
		w.Write(bytes)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(dat))
}
