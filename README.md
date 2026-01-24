# Chirpy

A Twitter-like social media microservice built in Go with JWT authentication and PostgreSQL.

## Features

- User registration and authentication with JWT tokens
- Create, read, and delete chirps (posts) with profanity filtering
- Premium "Chirpy Red" user tier
- RESTful API with comprehensive endpoints
- Metrics and health monitoring
- Polka webhook integration

## Tech Stack

- **Go 1.25.5** - Backend language
- **PostgreSQL** - Database with lib/pq driver
- **JWT** - Authentication with refresh tokens
- **sqlc** - Type-safe SQL queries
- **argon2id** - Secure password hashing

## Quick Start

1. **Clone and setup:**
   ```bash
   git clone <repository-url>
   cd Chirpy
   go mod download
   ```

2. **Configure environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your database URL and secrets
   ```

3. **Run database migrations:**
   ```bash
   # Apply schema from sql/schema/
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`.

## API Endpoints

### User Management
- `POST /api/users` - Create user
- `PUT /api/users` - Update user
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh token
- `POST /api/revoke` - Revoke token

### Chirps
- `POST /api/chirps` - Create chirp
- `GET /api/chirps` - Get all chirps
- `GET /api/chirps/{chirpID}` - Get specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete chirp

### System
- `GET /api/healthz` - Health check
- `GET /admin/metrics` - Application metrics
- `POST /admin/reset` - Reset application state
- `POST /api/polka/webhooks` - Polka webhook

## Development

```bash
# Run tests
go test ./...

# Generate database code
sqlc generate

# Build
go build -o Chirpy

# Format code
go fmt ./...
```

## Environment Variables

Required in `.env` file:
- `DB_URL` - PostgreSQL connection string
- `PLATFORM` - Platform identifier
- `SECRET_KEY` - JWT signing secret
- `POLKA_KEY` - Polka webhook key

## Contributing

1. Follow Go conventions and the patterns established in the codebase
2. Add tests for new features
3. Run `go fmt` and tests before submitting
