package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sianwa11/my-journal/internal/database"
)

func (cgf *apiConfig) createProject(w http.ResponseWriter, r *http.Request) {
	type Tags struct {
		ID   int `json:"id"`
		Name string `json:"name"`
	}

	type Params struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ImageUrl    string `json:"image_url"`
		Link        string `json:"link"`
		Github      string `json:"github"`
		Status      string `json:"status"`
		Tags 				[]Tags `json:"tags"`
	}

	var params Params
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	if params.Title == "" || params.Description == "" || params.Link == "" || params.Github == "" {
		respondWithError(w, http.StatusBadRequest, "please fill in required fields", nil)
		return
	}

	userIDInt := r.Context().Value(userIDKey).(int)
	userID := int64(userIDInt)

	tx, err := cgf.dbConn.Begin()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to start transaction", err)
		return
	}
	defer tx.Rollback()

	qtx := cgf.DB.WithTx(tx)

	project, err := qtx.CreateProject(r.Context(), database.CreateProjectParams{
		Title: params.Title,
		Description: params.Description,
		ImageUrl: sql.NullString{String: params.ImageUrl, Valid: params.ImageUrl != ""},
		Link: sql.NullString{String: params.Link, Valid: params.Link != ""},
		Github: sql.NullString{String: params.Github, Valid: params.Github != ""},
		Status: sql.NullString{String: params.Status, Valid: params.Status != ""},
		UserID: int64(userID),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create project", err)
		return
	}

	for _, tag := range params.Tags {
		if tag.ID == 0 && tag.Name ==  "" {
			respondWithError(w, http.StatusBadRequest, "invalid tag", err)
			return
		}

		var tagID int64
    // If tag has an ID, use it directly
    if tag.ID != 0 {
			tagID = int64(tag.ID)
    } else {
			// Try to find existing tag by name
			existingTag, err := qtx.SelectTag(r.Context(), tag.Name)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					// Tag doesn't exist, create it
					newTag, err := qtx.CreateTag(r.Context(), tag.Name)
					if err != nil {
							respondWithError(w, http.StatusInternalServerError, "failed to create tag", err)
							return
					}
					tagID = newTag.ID
				} else {
					// Some other database error
					respondWithError(w, http.StatusInternalServerError, "something went wrong getting tags", err)
					return
				}
			} else {
				// Tag exists, use its ID
				tagID = existingTag.ID
			}
    }


		_, err = qtx.CreateProjectTag(r.Context(), database.CreateProjectTagParams{
			ProjectID: project.ID,
			TagID: tagID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to create project tag", err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to commit transaction", err)
		return
	}

	respondWithJson(w, http.StatusCreated, struct{
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ImageUrl    string `json:"image_url"`
		Link        string `json:"link"`
		Github      string `json:"github"`
		Status      string `json:"status"`
		UserID			int    `json:"user_id"`
		Tags 				[]Tags `json:"tags"`
	}{
		ID: int(project.ID),
		Title: project.Title,
		Description: project.Description,
		ImageUrl: project.ImageUrl.String,
		Link: project.Link.String,
		Github: project.Github.String,
		Status: project.Status.String,
		UserID: int(userID),
		Tags: params.Tags,
	})
}