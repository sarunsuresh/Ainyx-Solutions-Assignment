-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, role, dob)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, email, role, dob, created_at;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, role, dob
FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT id, name, email, role, dob
FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users SET name = $1, email = $2, dob = $3
WHERE id = $4
RETURNING id, name, email, dob;

-- name: UpdatePassword :one
UPDATE users SET password_hash = $1
WHERE id = $2
RETURNING id, name, email, dob;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT id, name, dob FROM users ORDER BY id LIMIT $1 OFFSET $2;
