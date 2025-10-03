package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/sianwa11/my-journal/internal/database"
)

type Project struct {
	ProjectID   int       `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Link        string    `json:"link"`
	Github      string    `json:"github"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      int       `json:"user_id"`
	Tags        string    `json:"tags"`
}

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

func (cfg *apiConfig) getProjects(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}

	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid limit", err)
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid offset", err)
		return
	}

	projects, err := cfg.DB.GetProjects(r.Context(), database.GetProjectsParams{
		Limit: int64(limitInt),
		Offset: int64(offsetInt),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get projects", err)
		return
	}

	projectsArr := []Project{}
	for _, project := range projects {
		projectsArr = append(projectsArr, Project{
			ProjectID: int(project.ID),
			Title: project.Title,
			Description: project.Description,
			ImageURL: project.ImageUrl.String,
			Link: project.Link.String,
			Github: project.Github.String,
			Status: project.Status.String,
			CreatedAt: project.CreatedAt.Time,
			UpdatedAt: project.UpdatedAt.Time,
			UserID: int(project.UserID),
		})
	}

	respondWithJson(w, http.StatusOK, projectsArr)
}

func (cfg *apiConfig) getProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("projectID")
	if projectIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "empty projectID", nil)
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadGateway, "invalid projectID", err)
		return
	}

	project, err := cfg.DB.GetProject(r.Context(), int64(projectID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "project not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to get project", err)
		return
	}

	respondWithJson(w, http.StatusOK, Project{
		ProjectID: int(project.ProjectID),
		Title: project.Title,
		Description: project.Description,
		ImageURL: project.ImageUrl.String,
		Link: project.Link.String,
		Github: project.Github.String,
		Status: project.Status.String,
		CreatedAt: project.CreatedAt.Time,
		UserID: int(project.UserID),
		Tags: project.Tags,
	})
}

func (cfg *apiConfig) deleteProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("projectID")
	if projectIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing projectID", nil)
		return
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid projectID", err)
	}

	err = cfg.DB.DeleteProject(r.Context(), int64(projectID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete project", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}