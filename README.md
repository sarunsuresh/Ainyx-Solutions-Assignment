# User API

A production-style REST API built with **Go**, **Fiber**, and **PostgreSQL** demonstrating layered architecture, authentication, WebSockets, Redis, gRPC microservices, background workers, circuit breakers, and resilient email delivery.

---

# Features

- User CRUD
- JWT Authentication
- Role-based Authorization (User/Admin)
- Account Activation via Email
- gRPC Email Microservice
- WebSocket Notifications
- Redis-backed Rate Limiting
- Circuit Breaker
- Retry Queue with Exponential Backoff
- Queue Worker
- Global Error Handler
- Panic Recovery Middleware
- Structured Zap Logging
- Request ID Tracking
- Health Checks
- SQLC
- PostgreSQL

---

# Tech Stack

- **Go**
- **Fiber**
- **PostgreSQL**
- **Redis**
- **SQLC**
- **Zap**
- **JWT**
- **gRPC**
- **Protocol Buffers**
- **Docker**
- **Validator**

---

# Architecture

```
                        HTTP Client
                             │
                             ▼
                        Fiber Router
                             │
                 ┌───────────┴───────────┐
                 │                       │
           Global Middleware      WebSocket Endpoint
                 │                       │
                 ▼                       ▼
             Handlers              Hub Broadcast
                 │
                 ▼
             Services
                 │
      ┌──────────┴──────────┐
      │                     │
Repositories          Activation Service
      │                     │
      │                     ▼
      │             Email gRPC Client
      │                     │
      │              Circuit Breaker
      │                     │
      │              gRPC Email Service
      │
      ▼
 PostgreSQL

Background Components

Queue
 ↓
Worker
 ↓
Retry with Exponential Backoff
```

---

# Project Structure

```
cmd/
│
├── server/
│     main.go
│
└── emailservice/
      main.go

config/

db/
├── migrations/
└── sqlc/

internal/

├── activation/
│      account activation logic
│
├── clients/
│      gRPC client
│      circuit breaker
│
├── email/
│      gRPC email server
│
├── handler/
│
├── logger/
│
├── middleware/
│
├── models/
│
├── queue/
│      retry queue
│      worker
│
├── redis/
│
├── repository/
│
├── routes/
│
├── service/
│
└── websocket/
```

---

# Request Flow

```
HTTP Request

↓

Routes

↓

Middleware

↓

Handler

↓

Service

↓

Repository

↓

PostgreSQL
```

Email activation flow

```
Signup

↓

Activation Service

↓

Email Client

↓

Circuit Breaker

↓

gRPC Email Service
```

Retry flow

```
Email Failure

↓

Queue

↓

Worker

↓

Retry

↓

Email Service
```

---

# Authentication

JWT-based authentication.

Protected endpoints require

```
Authorization: Bearer <token>
```

Supports

- User
- Admin

Role-based authorization middleware protects admin endpoints.

---

# API Endpoints

## Authentication

| Method | Endpoint |
|---------|----------|
| POST | /auth/signup |
| POST | /auth/login |
| POST | /auth/activate |
| POST | /auth/resend |

---

## Users

| Method | Endpoint |
|---------|----------|
| GET | /users |
| GET | /users/me |
| GET | /users/:id |
| PUT | /users/:id |
| PUT | /users/:id/profile |
| PATCH | /users/me/passwordupdate |

---

## Admin

| Method | Endpoint |
|---------|----------|
| GET | /admin/users |
| DELETE | /admin/users/:id |

---

## Health

| Method | Endpoint |
|---------|----------|
| GET | /healthz |

---

## WebSocket

```
ws://localhost:3000/ws
```

Broadcasts

- user_created
- user_updated
- user_deleted
- user_profile_updated

---

# Email Service

Runs as a separate process.

```
cmd/emailservice
```

Responsibilities

- Receive gRPC requests
- Simulate email sending
- Artificial failures
- Health checks
- Structured logging

---

# Circuit Breaker

Protects the User API from repeatedly calling an unhealthy Email Service.

States

```
Closed

↓

Open

↓

Half Open

↓

Closed
```

Features

- Configurable failure threshold
- Fast-fail when open
- Automatic recovery testing
- Timeout-based reset

---

# Retry Queue

Failed emails are queued for retry.

Features

- Thread-safe queue
- Background worker
- Exponential backoff
- Automatic retries

Backoff

```
1 min

↓

2 min

↓

4 min

↓

8 min

↓

16 min

↓

30 min (cap)
```

---

# Logging

Zap structured logging records

- Request ID
- User ID
- Method
- Path
- Status
- Latency
- Client IP
- Redis failures
- WebSocket events
- Queue activity
- Circuit breaker state changes
- Transaction logs
- Panic stack traces

---

# Middleware

Implemented middleware

- Request ID
- Request Logger
- Recovery
- Global Error Handler
- JWT Authentication
- Role Authorization
- Redis Rate Limiter

---

# Error Handling

Global error handler converts application errors into consistent JSON responses.

Example

```json
{
    "error":"resource not found",
    "requestId":"..."
}
```

Handles

- Validation errors
- SQL errors
- Authentication errors
- Authorization errors
- Rate limits
- Panics
- Missing routes

---

# Health Checks

```
GET /healthz
```

Checks

- PostgreSQL
- Redis

Example

```json
{
    "status":"ok"
}
```

or

```json
{
    "status":"degraded",
    "redis":"down"
}
```

---

# Running

## Start everything

```
docker compose up --build
```

---

## Start only PostgreSQL + Redis

```
docker compose up db redis
```

Run User API

```
go run ./cmd/server
```

Run Email Service

```
go run ./cmd/emailservice
```

---

# Running Tests

```
go test ./...
```

Includes

- Service unit tests
- Circuit breaker tests
- Queue tests

---

# Technologies Demonstrated

- Layered Architecture
- Repository Pattern
- Dependency Injection
- JWT Authentication
- Role-Based Access Control
- SQLC
- PostgreSQL
- Redis
- WebSockets
- gRPC
- Protocol Buffers
- Circuit Breaker Pattern
- Retry Queue
- Background Workers
- Exponential Backoff
- Structured Logging
- Graceful Shutdown
- Docker
- Health Checks