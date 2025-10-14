package routes

import (
	"context"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sianwa11/my-journal/internal/database"
)

type apiConfig struct {
	dbConn    *sql.DB
	DB        *database.Queries
	jwtSecret string
}

func SetupRoutes() *http.ServeMux {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	secret := os.Getenv("SECRET")
	dbUrl := os.Getenv("DB_URL")
	// const rootPath = "."

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

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"split": strings.Split,
	}

	// Parse templates
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("template/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("template/partials/*.html"))

	getUserTemplateData := func() (map[string]interface{}, error) {
		user, err := apiCfg.DB.ListUser(context.Background())
		if err != nil || len(user) == 0 {
			return map[string]interface{}{
				"Name": "Your Name",
				"Year": time.Now().Year(),
			}, nil
		}

		return map[string]interface{}{
			"Name":     user[0].Name,
			"Email":    user[0].Email.String,
			"Github":   user[0].Github.String,
			"Linkedin": user[0].Linkedin.String,
			"Year":     time.Now().Year(),
		}, nil
	}

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

	mux.HandleFunc("/admin/profile", func(w http.ResponseWriter, r *http.Request) {
		user, err := apiCfg.DB.ListUser(r.Context())
		if err != nil {
			http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "profile.html", map[string]interface{}{
			"Title":    "My Profile",
			"Name":     user[0].Name,
			"Email":    user[0].Email.String,
			"Bio":      user[0].Bio.String,
			"Github":   user[0].Github.String,
			"Linkedin": user[0].Linkedin.String,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user, err := apiCfg.DB.ListUser(r.Context())
		if err != nil {
			http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
			return
		}

		if len(user) == 0 {
			http.Error(w, "No user data found", http.StatusInternalServerError)
			return
		}

		bio := ""
		if user[0].Bio.Valid {
			bio = strings.TrimSpace(user[0].Bio.String)
		}

		err = tmpl.ExecuteTemplate(w, "me.html", map[string]interface{}{
			"Title":       "Sianwa",
			"Name":        user[0].Name,
			"Bio":         bio,
			"Github":      user[0].Github.String,
			"Linkedin":    user[0].Linkedin.String,
			"Email":       user[0].Email.String,
			"CurrentPage": "about",
			"FooterText":  "Built with passion ❤️.",
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		data, _ := getUserTemplateData()
		data["Title"] = "Projects"
		data["CurrentPage"] = "projects"
		data["FooterText"] = "Built with passion ❤️."

		err := tmpl.ExecuteTemplate(w, "list-projects.html", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/journals", func(w http.ResponseWriter, r *http.Request) {
		data, _ := getUserTemplateData()
		data["Title"] = "Journals"
		data["CurrentPage"] = "journals"
		data["FooterText"] = "Built with passion ❤️."

		err := tmpl.ExecuteTemplate(w, "list-journals.html", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/journals/{ID}", func(w http.ResponseWriter, r *http.Request) {
		data, _ := getUserTemplateData()
		data["Title"] = "Journal Entry"
		data["CurrentPage"] = "journals"
		data["FooterText"] = "Built with passion ❤️."

		IDStr := r.PathValue("ID")
		journalID, err := strconv.Atoi(IDStr)

		if err != nil {
			http.Error(w, "Invalid journal ID", http.StatusBadRequest)
			return
		}

		journal, err := apiCfg.DB.GetJournalEntry(r.Context(), int64(journalID))
		if err != nil {
			http.Error(w, "Journal entry not found", http.StatusNotFound)
			return
		}

		nextAndPrev, err := apiCfg.DB.GetPrevAndNextJournalIDs(r.Context(), int64(journalID))
		if err != nil {
			http.Error(w, "Failed to fetch navigation data", http.StatusInternalServerError)
			return
		}

		data["Journal"] = journal
		data["NextJournalID"] = nextAndPrev.NextID
		data["PrevJournalID"] = nextAndPrev.PreviousID

		err = tmpl.ExecuteTemplate(w, "view-journal.html", data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/projects/{ID}", func(w http.ResponseWriter, r *http.Request) {
		data, _ := getUserTemplateData()
		data["Title"] = "Project Details"
		data["CurrentPage"] = "projects"
		data["FooterText"] = "Built with passion ❤️."

		IDStr := r.PathValue("ID")
		projectID, err := strconv.Atoi(IDStr)

		if err != nil {
			http.Error(w, "Invalid project ID", http.StatusBadRequest)
			return
		}

		project, err := apiCfg.DB.GetProject(r.Context(), int64(projectID))
		if err != nil {
			http.Error(w, "Project not found", http.StatusNotFound)
		}

		ordered, err := apiCfg.DB.GetProjectsNextAndPrevious(r.Context(), int64(projectID))
		if err != nil {
			http.Error(w, "Failed to fetch navigation data", http.StatusInternalServerError)
			return
		}

		data["Project"] = project
		data["NextProjectID"] = ordered.NextID
		data["PrevProjectID"] = ordered.PreviousID

		err = tmpl.ExecuteTemplate(w, "view-project.html", data)

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

	mux.HandleFunc("PUT /api/me", apiCfg.middlewareMustBeLoggedIn(apiCfg.editUserInfo))

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevokeToken)

	return mux
}
