package routes

import (
	"encoding/json"
	"fmt"
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

	journal, err := cfg.DB.CreateJournalEntry(r.Context(), database.CreateJournalEntryParams{Title: req.Title, Content: req.Content, UserID: 1})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create journal entry", err)
	}

	fmt.Printf("journal: %v\n", journal)
}