package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sianwa11/my-journal/internal/auth"
	"github.com/sianwa11/my-journal/internal/database"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var params Req
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode body", err)
		return
	}

	if params.Name == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Name and password are required", err)
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


	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to hash password", err)
		return
	}
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name: params.Name,
		Password: password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create user", err)
		return
	}

	respondWithJson(w, http.StatusCreated, struct {
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
	}{
		Name: user.Name,
		CreatedAt: user.CreatedAt.Time.String(),
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
		respondWithError(w, http.StatusInternalServerError, "failed to decode body", err)
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

	

	jwt, err := auth.MakeJWT(int(user.ID), cfg.jwtSecret, 1 * time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create jwt", err)
		return
	}

	respondWithJson(w, http.StatusOK, struct {
		ID        int `json:"id"`
		Name      string `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		Token     string `json:"token"`
	}{
		ID: int(user.ID),
		Name: user.Name,
		CreatedAt: user.CreatedAt.Time,
		Token: jwt,
	})
}

