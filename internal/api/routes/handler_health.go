package routes

import "net/http"

func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, map[string]string{
	"status": "healthy",
	"service": "my-journal",
	})
}