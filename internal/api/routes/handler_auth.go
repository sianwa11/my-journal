package routes

import (
	"context"
	"net/http"

	"github.com/sianwa11/my-journal/internal/auth"
)

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

		ctx := context.WithValue(r.Context(), "user_id", userId)


		next.ServeHTTP(w, r.WithContext(ctx))
	})
}