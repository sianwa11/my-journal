package routes

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sianwa11/my-journal/internal/database"
)

type apiConfig struct{
	DB        *database.Queries
	jwtSecret string
}


func SetupRoutes() *http.ServeMux{

	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
	
	secret := os.Getenv("SECRET")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		panic("Failed to open database " + err.Error())
	}

	
	apiCfg := &apiConfig{
		jwtSecret: secret,
	}
	apiCfg.DB = database.New(db)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/api/healthz", healthCheck)

	mux.HandleFunc("GET /api/journals", apiCfg.getJournalEntries)
	mux.HandleFunc("GET /api/journals/{journalID}", apiCfg.getJournalEntry)
	mux.HandleFunc("POST /api/journals", apiCfg.middlewareMustBeLoggedIn(apiCfg.postJournalEntry))
	mux.HandleFunc("PUT /api/journals", apiCfg.middlewareMustBeLoggedIn(apiCfg.editJournalEntry))
	mux.HandleFunc("DELETE /api/journals/{journalID}", apiCfg.middlewareMustBeLoggedIn(apiCfg.deleteJournalEntry))

	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevokeToken)


	return mux
}