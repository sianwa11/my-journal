package routes

import (
	"encoding/json"
	"net/http"

	"github.com/sianwa11/my-journal/internal/database"
)

func (cfg *apiConfig) postJournalEntry(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var req Req
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode body", err)
		return
	}

	userIdInt := r.Context().Value("user_id").(int)
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