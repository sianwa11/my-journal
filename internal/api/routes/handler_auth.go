package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sianwa11/my-journal/internal/auth"
	"github.com/sianwa11/my-journal/internal/database"
)

type contextKey string

const userIDKey contextKey = "user_id"

func (cfg *apiConfig) middlewareMustBeLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "missing or invalid authorization header", err)
			return
		}

		userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid or expired token", err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var params Params
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON format", err)
		return
	}

	if params.Name == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Name and Password required", err)
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get user", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Incorrect password", err)
		return
	}

	jwt, err := auth.MakeJWT(int(user.ID), cfg.jwtSecret, 1*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create jwt", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create refresh token", err)
		return
	}

	refreshTokenDB, err := cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    int64(user.ID),
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to save refresh token", err)
		return
	}

	respondWithJson(w, http.StatusOK, struct {
		ID           int       `json:"id"`
		Name         string    `json:"name"`
		CreatedAt    time.Time `json:"created_at"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}{
		ID:           int(user.ID),
		Name:         user.Name,
		CreatedAt:    user.CreatedAt.Time,
		Token:        jwt,
		RefreshToken: refreshTokenDB.Token,
	})
}

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid authorization header", err)
		return
	}

	userId, err := cfg.DB.GetByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(int(userId), cfg.jwtSecret, 1*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create jwt", err)
		return
	}

	respondWithJson(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid refresh token", err)
		return
	}

	err = cfg.DB.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
