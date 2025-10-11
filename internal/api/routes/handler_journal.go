package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sianwa11/my-journal/internal/database"
)

type Journal struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserID    int    `json:"user_id"`
}

type JournalsResponse struct {
	Journals []Journal `json:"journals"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
	HasMore  bool      `json:"has_more"`
}

func (cfg *apiConfig) getJournalEntries(w http.ResponseWriter, r *http.Request) {

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
		respondWithError(w, http.StatusBadRequest, "invalid limit parameter", err)
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid offset parameter", err)
		return
	}

	totalCount, err := cfg.DB.GetAllJournalsCount(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get journals count", err)
		return
	}

	journals, err := cfg.DB.GetJournals(r.Context(), database.GetJournalsParams{
		Limit: int64(limitInt),
		Offset: int64(offsetInt),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not fetch journals", err)
		return
	}


	journalEntries := []Journal{}
	for _, journal := range journals {
		journalEntries = append(journalEntries, Journal{
			ID: int(journal.ID),
			Title: journal.Title,
			Content: journal.Content,
			CreatedAt: journal.CreatedAt.Time.String(),
			UpdatedAt: journal.UpdatedAt.Time.String(),
			UserID: int(journal.UserID),
		})
	}

		// Calculate pagination info
  currentPage := (offsetInt / limitInt) + 1
  hasMore := offsetInt + len(journals) < int(totalCount)

	respondWithJson(w, http.StatusOK, JournalsResponse {
		Journals: journalEntries,
		Total: int(totalCount),
		Page: currentPage,
		Limit: limitInt,
		HasMore: hasMore,
	})
}

func (cfg *apiConfig) getJournalEntry(w http.ResponseWriter, r *http.Request) {
	journalIDString := r.PathValue("journalID")
	journalID, err := strconv.Atoi(journalIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid journal ID", err)
		return
	}

	journalEntry, err := cfg.DB.GetJournalEntry(r.Context(), int64(journalID))
	if err != nil {
		if errors.Is(err,sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "journal not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to get journal", err)
		return
	}

	respondWithJson(w, http.StatusOK, Journal{
		ID: int(journalEntry.ID),
		Title: journalEntry.Title,
		Content: journalEntry.Content,
		CreatedAt: journalEntry.CreatedAt.Time.String(),
		UpdatedAt: journalEntry.CreatedAt.Time.String(),
		UserID: int(journalEntry.UserID),
	})

}

func (cfg *apiConfig) postJournalEntry(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req Req
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON format", err)
		return
	}

	userIdInt := r.Context().Value(userIDKey).(int)
	userId := int64(userIdInt)

	journal, err := cfg.DB.CreateJournalEntry(r.Context(), database.CreateJournalEntryParams{
		Title: req.Title,
		Content: req.Content,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create journal entry", err)
	}

	respondWithJson(w, http.StatusCreated, struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		UserID    int    `json:"user_id"`
	}{
		 Title: journal.Title,
		 Content: journal.Content,
		 CreatedAt: journal.CreatedAt.Time.String(),
		 UserID: int(journal.UserID),
	})

}

func (cfg *apiConfig) editJournalEntry(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var params Params
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON format", err)
		return
	}

	if params.Title == "" || params.ID == 0 || params.Content == "" {
		respondWithError(w, http.StatusBadRequest, "please fill in required fields",err)
		return
	}

	err = cfg.DB.UpdateJournalEntry(r.Context(), database.UpdateJournalEntryParams{
		Title: params.Title,
		Content: params.Content,
		ID: int64(params.ID),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to update journal", err)
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "journal updated successfully",
	})	
}

func (cfg *apiConfig) deleteJournalEntry(w http.ResponseWriter, r *http.Request) {
	journalIDString := r.PathValue("journalID")
	journalID, err := strconv.Atoi(journalIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid journal ID", err)
		return
	}

	userID := r.Context().Value(userIDKey).(int)

	journal, err := cfg.DB.GetUsersJournal(r.Context(), database.GetUsersJournalParams{
		ID: int64(journalID),
		UserID: int64(userID),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	if journal.UserID != int64(userID) {
		respondWithError(w, http.StatusUnauthorized, "cannot perform this action", err)
		return
	}

	err = cfg.DB.DeleteJournalEntry(r.Context(), database.DeleteJournalEntryParams{
		ID: int64(journalID),
		UserID: int64(userID),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete journal", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}