package routes

import (
	"encoding/json"
	"net/http"

	"github.com/sianwa11/my-journal/internal/database"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Name string `json:"name"`
	}

	var params Req
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode body", err)
		return
	}

	if params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Name is required", err)
		return
	}

	// check if users exists
	users, err := cfg.DB.ListUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get users", err)
		return
	}

	if len(users) > 0 {
		respondWithError(w, http.StatusForbidden, "another account already exists", err)
		return
	}


	apiKey := "random-key"
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name: params.Name,
		ApiKey: apiKey,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create user", err)
		return
	}

	respondWithJson(w, http.StatusCreated, struct {
		Name   string `json:"name"`
		ApiKey string `json:"api_key"`
	}{
		Name: user.Name,
		ApiKey: user.ApiKey,
	})

}