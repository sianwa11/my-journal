package routes

import (
	"net/http"
)

type Tag struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

func (cfg *apiConfig) searchTags(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		respondWithError(w, http.StatusBadRequest, "missing search query", nil)
		return
	}

	tags, err := cfg.DB.SearchTags(r.Context(), "%"+query+"%")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to search tags", err)
		return
	}

	tagsArr := []Tag{}
	for _, tag := range tags {
		tagsArr = append(tagsArr, Tag{
			ID:    int(tag.ID),
			Value: tag.Value,
		})
	}

	respondWithJson(w, http.StatusOK, tagsArr)
}
