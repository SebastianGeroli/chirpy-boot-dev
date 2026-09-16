package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SebastianGeroli/chirpy-boot-dev/internal/auth"
	"github.com/SebastianGeroli/chirpy-boot-dev/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	userParams := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          params.Email,
		HashedPassword: hashed_password,
	}

	createdUser, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, 500, err.Error())
	}
	user := User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}
	respondWithJSON(w, 201, user)
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPassword(params.Password, dbUser.HashedPassword)
	if err != nil || !match {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJSON(w, 200, user)
}
