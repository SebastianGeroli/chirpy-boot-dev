package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/SebastianGeroli/chirpy-boot-dev/internal/auth"
	"github.com/SebastianGeroli/chirpy-boot-dev/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	requestedID := r.PathValue("chirpID")
	ID, err := uuid.Parse(requestedID)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}
	dbChirp, err := cfg.db.GetChirp(r.Context(), ID)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJSON(w, 200, chirp)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		newChirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
		chirps = append(chirps, newChirp)
	}
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)

	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	chirpMsg, err := validate_chirp(params.Body)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	chirpParams := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      chirpMsg,
		UserID:    userID,
	}
	chirpCreated, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	chirp := Chirp{
		ID:        chirpCreated.ID,
		CreatedAt: chirpCreated.CreatedAt,
		UpdatedAt: chirpCreated.UpdatedAt,
		Body:      chirpCreated.Body,
		UserID:    chirpCreated.UserID,
	}
	respondWithJSON(w, 201, chirp)
}

func validate_chirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		return "", errors.New("Chirp is too long")
	}

	words := strings.Split(chirp, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if loweredWord == "kerfuffle" || loweredWord == "sharbert" || loweredWord == "fornax" {
			words[i] = "****"
		}
	}

	sanitazedChirp := strings.Join(words, " ")
	return sanitazedChirp, nil
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, 403, "Not author of chirp")
		return
	}

	_, err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}

	respondWithJSON(w, 204, nil)
}
