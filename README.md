# Ticket System - Backend Intern Assignment

A REST API ticket system built with **Golang**, **Gin**, **GORM**, and **SQLite**.

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin
- **ORM**: GORM with SQLite
- **Auth**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt

## Local Run (without Docker)

```bash
go mod tidy
go run ./cmd/main.go
```

## Docker Run

```bash
docker build -t ticket-system .
docker run -p 8080:8080 ticket-system
```

Health check:
```bash
curl http://localhost:8080/health
```

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /health | No | Health check |
| POST | /auth/register | No | Register user |
| POST | /auth/login | No | Login, returns JWT |
| POST | /tickets | Yes | Create ticket |
| GET | /tickets | Yes | List own tickets |
| GET | /tickets/:id | Yes | Get own ticket by ID |
| PATCH | /tickets/:id/status | Yes | Update ticket status |

## Example Usage

### Register
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Create Ticket
```bash
curl -X POST http://localhost:8080/tickets \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Fix login bug","detail":"Login fails on mobile"}'
```

### Update Status
```bash
curl -X PATCH http://localhost:8080/tickets/<id>/status \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'
```

## Status Flow

```
open → in_progress → closed
```

- A closed ticket cannot be reopened.

## Deployment

Deployed on Render (free tier):  
**URL**: `https://YOUR-APP.onrender.com`  
**Health Check**: `https://YOUR-APP.onrender.com/health`

## Assumptions

- SQLite is used for simplicity and zero-config deployment.
- JWT expiry is set to 24 hours.
- A user can only view and update their own tickets.
