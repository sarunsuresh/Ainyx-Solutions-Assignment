# User API

A simple REST API built with Go and Fiber to manage users. Each user has a name and a date of birth, and the API calculates their age automatically when you fetch them.

## Tech Stack

- **Go** with **Fiber** (web framework)
- **PostgreSQL** (database)
- **SQLC** (generates type-safe Go code from SQL)
- **Zap** (logging)
- **validator** (input validation)

## Project Structure

```
cmd/server/main.go     # starts the app
config/                # reads settings from environment
db/migrations/         # database table setup
db/sqlc/               # SQL queries + generated code
internal/
  handler/             # handles HTTP requests/responses
  service/             # business logic (age calculation)
  repository/          # database access
  routes/              # connects URLs to handlers
  middleware/          # request ID + request logging
  models/              # request/response shapes
  logger/              # logger setup
```

Request flow: **routes → handler → service → repository → database**

## Getting Started

### Option 1: Run everything with Docker

```bash
docker-compose up --build
```

The database and app both start, and the table is created automatically.

### Option 2: Run the database in Docker, app in your terminal

Start just the database:
```bash
docker-compose up db
```

In a second terminal, run the app (pointing it at the Docker database on port 5433):
```bash
# Windows PowerShell
$env:DB_HOST="localhost"
$env:DB_PORT="5433"
go run ./cmd/server
```

The server runs at `http://localhost:3000`.

## API Endpoints

| Method | Endpoint      | Description              |
|--------|---------------|--------------------------|
| POST   | `/users`      | Create a user            |
| GET    | `/users`      | List all users           |
| GET    | `/users/:id`  | Get one user (with age)  |
| PUT    | `/users/:id`  | Update a user            |
| DELETE | `/users/:id`  | Delete a user            |

The list endpoint supports pagination: `/users?limit=10&offset=0`

### Example: Create a user

POST `/users`
```json
{
  "name": "Alice",
  "dob": "1990-05-10"
}
```

Response:
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Example: Get a user

GET `/users/1`
```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 36
}
```

> Note: `dob` must be in `YYYY-MM-DD` format. Age is calculated automatically.

## Running Tests

```bash
go test ./...
```

The age calculation has unit tests in `internal/service/user_service_test.go`.

## Notes

- Age is calculated on every request instead of being stored, so it's always accurate.
- Every response includes an `X-Request-ID` header for tracking.
- Each request is logged with its method, path, status, and duration.
