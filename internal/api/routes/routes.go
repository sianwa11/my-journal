package routes

import (
	"net/http"
)

type apiConfig struct{}


func SetupRoutes() *http.ServeMux{
	mux := http.NewServeMux()

	apiCfg := &apiConfig{}

	mux.HandleFunc("/api/healthz", healthCheck)

	mux.HandleFunc("POST /api/journal", apiCfg.postJournalEntry)

	return mux
}