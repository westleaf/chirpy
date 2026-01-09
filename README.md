# Chirpy

A sample implementation social media API server built with Go that allows users to post short messages called "chirps", similar to Twitter.

## What It Does

Chirpy is a RESTful API that provides:

- User registration and authentication with JWT tokens
- Creating, reading, and deleting chirps (short messages)
- User profile management
- Refresh token system for secure authentication
- Premium user upgrades (Chirpy Red)

## Why Use This Project

- Clean architecture with separate database and authentication layers
- Production-ready patterns including middleware and proper error handling
- Secure password hashing with Argon2id
- PostgreSQL database with SQL migrations
- Type-safe database queries using sqlc
- Environment-based configuration

## Prerequisites

- Go 1.25.4 or higher
- PostgreSQL
- Goose (database migration tool)
- sqlc (SQL code generator)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/westleaf/chirpy.git
cd chirpy
```

2. Install dependencies:
```bash
go mod download
```

3. Install required tools:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

4. Set up PostgreSQL database:
```bash
createdb chirpy
```

5. Create a `.env` file in the project root with the following variables:
```
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWTSECRET=your_jwt_secret_here
POLKA_KEY=your_polka_api_key_here
```

6. Run database migrations:
```bash
make up
```

7. Generate database code:
```bash
make sqlc
```

## Running the Server

Build and run the server:
```bash
go build -o chirpy && ./chirpy
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health & Metrics
- `GET /api/healthz` - Health check endpoint
- `GET /admin/metrics` - View server metrics
- `POST /admin/reset` - Reset database (dev only)

### Users
- `POST /api/users` - Create a new user
- `PUT /api/users` - Update user email/password
- `POST /api/login` - Login and receive JWT tokens

### Chirps
- `GET /api/chirps` - Get all chirps
- `GET /api/chirps/{id}` - Get a specific chirp
- `POST /api/chirps` - Create a new chirp (requires authentication)
- `DELETE /api/chirps/{id}` - Delete a chirp (requires authentication)

### Authentication
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token

## Database Management

Run migrations:
```bash
make up          # Run all migrations
make down        # Roll back last migration
make status      # Check migration status
```

Regenerate database code after SQL changes:
```bash
make sqlc
```

## Project Structure

```
.
├── main.go              # Application entry point
├── internal/            # Private application code
│   ├── auth/           # Authentication logic
│   └── database/       # Generated database queries
├── sql/
│   ├── queries/        # SQL query definitions
│   └── schema/         # Database migrations
└── assets/             # Static files
```

## License

This project is provided as-is for educational purposes. Based on the boot.dev course Learn HTTP servers in Go.
