package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/sianwa11/my-journal/internal/database"
)

func (cfg *apiConfig) editBio(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Bio string `json:"bio"`
	}

	var params Params
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	userID := r.Context().Value(userIDKey).(int)
	if userID == 0 {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	err = cfg.DB.UpdateBio(r.Context(), database.UpdateBioParams{
		Bio: sql.NullString{String: params.Bio, Valid: params.Bio != ""},
		ID: int64(userID),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update bio", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}