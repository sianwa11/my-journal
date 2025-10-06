package routes

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sianwa11/my-journal/internal/database"
)

type apiConfig struct{
	dbConn        *sql.DB
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
	const rootPath = "."

	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		panic("Failed to open database " + err.Error())
	}

	
	apiCfg := &apiConfig{
		jwtSecret: secret,
	}
	apiCfg.DB = database.New(db)
	apiCfg.dbConn = db
	
	mux := http.NewServeMux()

	// Serves static files from the "static" directory
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Parse templates
	tmpl := template.Must(template.ParseGlob("template/*.html"))


	// Dashboard route using template
	mux.HandleFunc("/admin/dashboard", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{
			"Title": "Admin Dashboard",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"Title": "Admin Login",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/admin/journals", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "journals.html", map[string]interface{}{
			"Title": "Manage Journals",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/admin/projects", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "projects.html", map[string]interface{}{
			"Title": "Manage Projects",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/api/healthz", healthCheck)

	mux.HandleFunc("GET /api/journals", apiCfg.getJournalEntries)
	mux.HandleFunc("GET /api/journals/{journalID}", apiCfg.getJournalEntry)
	mux.HandleFunc("POST /api/journals", apiCfg.middlewareMustBeLoggedIn(apiCfg.postJournalEntry))
	mux.HandleFunc("PUT /api/journals", apiCfg.middlewareMustBeLoggedIn(apiCfg.editJournalEntry))
	mux.HandleFunc("DELETE /api/journals/{journalID}", apiCfg.middlewareMustBeLoggedIn(apiCfg.deleteJournalEntry))

	mux.HandleFunc("POST /api/projects", apiCfg.middlewareMustBeLoggedIn(apiCfg.createProject))
	mux.HandleFunc("GET /api/projects", apiCfg.getProjects)
	mux.HandleFunc("GET /api/projects/{projectID}", apiCfg.getProject)
	mux.HandleFunc("DELETE /api/projects/{projectID}", apiCfg.middlewareMustBeLoggedIn(apiCfg.deleteProject))
	mux.HandleFunc("PUT /api/projects", apiCfg.middlewareMustBeLoggedIn(apiCfg.updateProject))


	mux.HandleFunc("GET /api/tags", apiCfg.searchTags)

	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevokeToken)


	return mux
}