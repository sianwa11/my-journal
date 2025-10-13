package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/sianwa11/my-journal/internal/database"
)

func (cfg *apiConfig) editUserInfo(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Github   string `json:"github"`
		Linkedin string `json:"linkedin"`
		Bio      string `json:"bio"`
	}

	var params Params
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	if params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "name is required", nil)
		return
	}

	userID := r.Context().Value(userIDKey).(int)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	err = cfg.DB.UpdateUserInfo(r.Context(), database.UpdateUserInfoParams{
		Name:     params.Name,
		Email:    sql.NullString{String: params.Email, Valid: params.Email != ""},
		Github:   sql.NullString{String: params.Github, Valid: params.Github != ""},
		Linkedin: sql.NullString{String: params.Linkedin, Valid: params.Linkedin != ""},
		Bio:      sql.NullString{String: params.Bio, Valid: params.Bio != ""},
		ID:       int64(userID),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update bio", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
