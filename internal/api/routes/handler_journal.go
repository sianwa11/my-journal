package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) postJournalEntry(w http.ResponseWriter, r *http.Request) {
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

	fmt.Printf("Received new journal entry: Title=%s, Content=%s\n", req.Title, req.Content)
}