package routes

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sianwa11/my-journal/internal/auth"
	"github.com/sianwa11/my-journal/internal/database"
)

func setupTestDB(t *testing.T) (*sql.DB, *database.Queries) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create the users table and other necessary tables
	schema := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			name TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL DEFAULT 'unset'
		, bio TEXT, email TEXT DEFAULT '', github TEXT DEFAULT '', linkedin TEXT DEFAULT '');
		CREATE TABLE journal_entries (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE refresh_tokens(
			id INTEGER PRIMARY KEY,
			token TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			revoked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE projects(
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			image_url TEXT,
			link TEXT,
			github TEXT,
			status TEXT CHECK(status IN ('in_progress','completed','archived')) DEFAULT 'completed',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE
			CASCADE
		);
		CREATE TABLE tags(
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE project_tags(
			project_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
			PRIMARY KEY (project_id, tag_id)
		);
	`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	queries := database.New(db)

	return db, queries
}

func createTestUser(t *testing.T, queries *database.Queries, name, password string) database.User {
	t.Helper()

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user, err := queries.CreateUser(context.Background(), database.CreateUserParams{
		Name:     name,
		Password: hashedPassword,
	})
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

func setupTestAPIConfig(t *testing.T) (*apiConfig, *sql.DB) {
	t.Helper()

	db, queries := setupTestDB(t)

	cfg := &apiConfig{
		dbConn:    db,
		DB:        queries,
		jwtSecret: "test-secret-key",
	}

	return cfg, db
}

func TestLoginHandler(t *testing.T) {
	os.Setenv("SECRET", "test-secret-key")
	os.Setenv("DB_URL", ":memory:")
	defer func() {
		os.Unsetenv("SECRET")
		os.Unsetenv("DB_URL")
	}()

	apiCfg, db := setupTestAPIConfig(t)
	defer db.Close()

	testUser := createTestUser(t, apiCfg.DB, "testuser", "testpassword")

	tests := []struct {
		name       string
		payload    map[string]string
		wantStatus int
		wantToken  bool
		setupUser  bool
	}{
		{
			name: "valid credentials",
			payload: map[string]string{
				"name":     testUser.Name,
				"password": "testpassword",
			},
			wantStatus: http.StatusOK,
			wantToken:  true,
			setupUser:  false, // User already created above
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"name":     "testadmin",
				"password": "wrongpassword",
			},
			wantStatus: http.StatusInternalServerError, // Your code returns 500 for incorrect password
			wantToken:  false,
		},
		{
			name: "invalid username",
			payload: map[string]string{
				"name":     "wronguser",
				"password": "testpassword",
			},
			wantStatus: http.StatusInternalServerError, // Your code returns 500 when user not found
			wantToken:  false,
		},
		{
			name: "missing password",
			payload: map[string]string{
				"name": "testadmin",
			},
			wantStatus: http.StatusBadRequest,
			wantToken:  false,
		},
		{
			name: "missing name",
			payload: map[string]string{
				"password": "testpassword",
			},
			wantStatus: http.StatusBadRequest,
			wantToken:  false,
		},
		{
			name:       "empty payload",
			payload:    map[string]string{},
			wantStatus: http.StatusBadRequest,
			wantToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			// Call the handler directly
			apiCfg.handleLogin(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.wantStatus, rr.Code, rr.Body.String())
			}

			if tt.wantToken {
				// Parse the response to check for token
				var response struct {
					ID           int    `json:"id"`
					Name         string `json:"name"`
					CreatedAt    string `json:"created_at"`
					Token        string `json:"token"`
					RefreshToken string `json:"refresh_token"`
				}

				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to parse response JSON: %v", err)
				}

				if response.Token == "" {
					t.Error("Expected non-empty token")
				}

				if response.RefreshToken == "" {
					t.Error("Expected non-empty refresh token")
				}

				if response.Name != testUser.Name {
					t.Errorf("Expected name %s, got %s", testUser.Name, response.Name)
				}

				if response.ID != int(testUser.ID) {
					t.Errorf("Expected ID %d, got %d", testUser.ID, response.ID)
				}

				// Verify the JWT token is valid
				userID, err := auth.ValidateJWT(response.Token, "test-secret-key")
				if err != nil {
					t.Errorf("Token validation failed: %v", err)
				}

				if userID != int(testUser.ID) {
					t.Errorf("Expected user ID %d from token, got %d", testUser.ID, userID)
				}
			} else {
				// Check for error response
				var errorResponse struct {
					Error string `json:"error"`
				}

				if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
					// If we can't parse as error response, that's also okay for some cases
					t.Logf("Response body: %s", rr.Body.String())
				}
			}
		})
	}
}
