# Chirpy

A RESTful API server built in Go as part of the [boot.dev](https://boot.dev) backend development course. Chirpy is a simple social media backend that supports users, chirps (short messages), authentication, and refresh tokens.

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Auth:** JWT (access tokens) + refresh tokens
- **Password hashing:** Argon2id
- **DB migrations:** Goose
- **Query generation:** SQLC

## Setup

### Prerequisites

- Go 1.22+
- PostgreSQL
- [Goose](https://github.com/pressly/goose)
- [SQLC](https://sqlc.dev)

### Installation

```bash
git clone https://github.com/borisfritz/chirpy
cd chirpy
go mod tidy
```

### Environment Variables

Create a `.env` file in the project root:

```env
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=your_generated_secret_here
```

Generate a JWT secret:
```bash
openssl rand -base64 64
```

### Database Setup

```bash
make migrate-up
```

## Running the Server

```bash
make run
```

Server starts on `http://localhost:8080`.

## API Endpoints

### Health

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/healthz` | Readiness check |

### Users

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/users` | No | Create a user |
| POST | `/api/login` | No | Login and receive tokens |
| PUT | `/api/users` | JWT | Update user |

### Chirps

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/chirps` | JWT | Create a chirp |
| GET | `/api/chirps` | No | Get all chirps |
| GET | `/api/chirps/{chirpID}` | No | Get a single chirp |
| DELETE | `/api/chirps/{chirpID}` | JWT | Delete a chirp |

Query parameters for `GET /api/chirps`:
- `author_id` — filter chirps by user ID
- `sort` — `asc` (default) or `desc`

### Tokens

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/refresh` | Refresh token | Get a new access token |
| POST | `/api/revoke` | Refresh token | Revoke a refresh token |

### Admin

| Method | Path | Description |
|--------|------|-------------|
| GET | `/admin/metrics` | View server hit count |
| POST | `/admin/reset` | Reset hits and users (dev only) |

## Authentication

Chirpy uses a two-token auth system:

- **Access token (JWT)** — short lived (1 hour), included in `Authorization: Bearer <token>` header
- **Refresh token** — long lived (60 days), used to get new access tokens without re-logging in

## Running Tests

```bash
make test
```

## Makefile Commands

```bash
make run            # run the server
make build          # build the binary
make test           # run all tests with coverage
make migrate-up     # apply migrations
make migrate-down   # roll back last migration
make migrate-reset  # reset all migrations
```
