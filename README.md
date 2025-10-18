# My Journal

A personal journaling application built with Go, featuring a clean web interface and secure user authentication. This application allows users to create and manage their personal journal entries in a private, single-user environment.

## Features

- **Single-user system** - Only one account can be created per instance
- **Secure authentication** - JWT-based authentication with password hashing
- **Clean web interface** - Built with Tailwind CSS for a modern look
- **Database flexibility** - Supports both SQLite and Turso (libSQL)
- **Dockerized deployment** - Easy deployment with Docker and Google Cloud Run
- **Database migrations** - Managed with Goose migration tool

## Tech Stack

- **Backend**: Go 1.24+
- **Database**: SQLite/Turso (libSQL)
- **Frontend**: HTML templates with Tailwind CSS
- **Authentication**: JWT tokens with bcrypt password hashing
- **Deployment**: Docker, Google Cloud Run
- **Migrations**: Goose

## Prerequisites

- Go 1.24 or higher
- Node.js 18+ (for Tailwind CSS)
- Docker (for containerized deployment)
- Google Cloud CLI (for cloud deployment)

## Local Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/sianwa11/my-journal.git
   cd my-journal
   ```

2. **Install dependencies**
   ```bash
   # Go dependencies
   go mod download
   
   # Node.js dependencies for Tailwind CSS
   npm install
   ```

3. **Set up environment variables**
   Create a `.env` file in the project root:
   ```env
   DATABASE_URL=./data/my-journal.db
   JWT_SECRET=your-super-secret-jwt-key
   PORT=8080
   ```

4. **Install Goose for database migrations**
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

5. **Run database migrations**
   ```bash
   ./scripts/migrateup.sh
   ```

6. **Build CSS**
   ```bash
   npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
   ```

7. **Run the application**
   ```bash
   go run main.go
   ```

8. **Access the application**
   Open your browser and navigate to `http://localhost:8080`

## Production Deployment

### Using Docker

1. **Build production assets**
   ```bash
   ./scripts/buildprod.sh
   ```

2. **Build Docker image**
   ```bash
   docker build -t my-journal .
   ```

3. **Run container**
   ```bash
   docker run -p 8080:8080 -e DATABASE_URL="your-database-url" -e JWT_SECRET="your-jwt-secret" my-journal
   ```

### Google Cloud Run

The project includes GitHub Actions for automated deployment to Google Cloud Run:

1. **Set up required secrets in GitHub:**
   - `GCP_CREDENTIALS`: Service account JSON with necessary permissions
   - `DB_URL`: Your database connection string

2. **Push to main branch:**
   ```bash
   git push origin main
   ```

3. **The CI/CD pipeline will automatically:**
   - Build the application and CSS
   - Create Docker image
   - Push to Google Artifact Registry
   - Run database migrations
   - Deploy to Cloud Run

## Database Options

### SQLite (Local Development)
```env
DATABASE_URL=./data/my-journal.db
```

### Turso (Production)
```env
DATABASE_URL=libsql://your-database.turso.io?authToken=your-token
```

## API Endpoints

- `POST /api/users` - Create user account (only if no users exist)
- `POST /api/login` - User authentication
- `GET /api/entries` - Get journal entries
- `POST /api/entries` - Create journal entry
- `PUT /api/entries/{id}` - Update journal entry
- `DELETE /api/entries/{id}` - Delete journal entry

## Project Structure

```
├── internal/
│   ├── api/           # HTTP handlers and routes
│   ├── auth/          # Authentication logic
│   ├── database/      # Database models and queries
│   └── sql/           # Database migrations
├── static/            # CSS and static assets
├── template/          # HTML templates
├── scripts/           # Build and deployment scripts
├── .github/workflows/ # CI/CD configuration
└── main.go           # Application entry point
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and ensure code quality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Security

- Passwords are hashed using bcrypt
- JWT tokens for session management
- Single-user restriction prevents unauthorized access
- HTTPS recommended for production deployment

## Support

If you encounter any issues or have questions, please open an issue on GitHub.