package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/SebastianGeroli/chirpy-boot-dev/internal/auth"
	"github.com/SebastianGeroli/chirpy-boot-dev/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	refresh_token := auth.MakeRefreshToken()
	refreshParams := database.CreateRefreshTokenParams{
		Token:     refresh_token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}
	dbToken, err := cfg.db.CreateRefreshToken(r.Context(), refreshParams)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	user := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: dbToken.Token,
	}
	respondWithJSON(w, 200, user)
}

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	dbToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	if dbToken.RevokedAt.Valid {
		respondWithError(w, 401, "Revoked token")
		return
	}

	if dbToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "Refresh token expired")
		return
	}

	newToken, err := auth.MakeJWT(dbToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	type result struct {
		Token string `json:"token"`
	}

	response := result{
		Token: newToken,
	}
	respondWithJSON(w, 200, response)
}

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	params := database.RevokeTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Token:     token,
	}
	_, err = cfg.db.RevokeToken(r.Context(), params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 204, nil)
}
