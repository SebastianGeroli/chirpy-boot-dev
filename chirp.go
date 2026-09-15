package main

import (
	"encoding/json"
	"net/http"
)

func chirp(responseWriter http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnError struct {
		Error string `json:"error"`
	}

	type returnValid struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		error := returnError{
			Error: "Failed to, parse body. Invalid JSON",
		}
		dat, internalErr := json.Marshal(error)
		if internalErr != nil {
			responseWriter.WriteHeader(500)
			bytes := []byte("Failed to parse error.. internal server error")
			responseWriter.Write(bytes)
		}
		responseWriter.WriteHeader(400)
		responseWriter.Write([]byte(dat))
	}

	if len(params.Body) > 140 {
		error := returnError{
			Error: "Chirp is too long",
		}
		dat, internalErr := json.Marshal(error)
		if internalErr != nil {
			responseWriter.WriteHeader(500)
			bytes := []byte("Failed to parse error.. internal server error")
			responseWriter.Write(bytes)
		}
		responseWriter.WriteHeader(400)
		responseWriter.Write([]byte(dat))
	}

	validResponse := returnValid{
		Valid: true,
	}
	dat, err := json.Marshal(validResponse)
	if err != nil {
		responseWriter.WriteHeader(500)
		bytes := []byte("Failed to parse valid response... internal server error")
		responseWriter.Write(bytes)
	}
	responseWriter.WriteHeader(200)
	responseWriter.Write([]byte(dat))
}
