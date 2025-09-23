package routes

import (
	"database/sql"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sianwa11/my-journal/internal/database"
)

type apiConfig struct{
	DB *database.Queries
}


func SetupRoutes() *http.ServeMux{
	
	// FIXME: Get from an env
	db, err := sql.Open("sqlite3", "./my-journal.db")
	if err != nil {
		panic("Failed to open database " + err.Error())
	}
	
	apiCfg := &apiConfig{}
	dbQueries := database.New(db)
	apiCfg.DB = dbQueries
	
	mux := http.NewServeMux()
	mux.HandleFunc("/api/healthz", healthCheck)

	mux.HandleFunc("POST /api/journal", apiCfg.postJournalEntry)

	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)

	return mux
}