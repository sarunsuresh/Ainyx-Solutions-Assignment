-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, role, dob)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, email, role, dob, created_at;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, role, dob,is_active
FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT id, name, email, role, dob,is_active
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

-- name: UpsertAddress :one
INSERT INTO addresses (user_id, line1, line2, city, state, postal_code, country)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id)
DO UPDATE SET
    line1       = EXCLUDED.line1,
    line2       = EXCLUDED.line2,
    city        = EXCLUDED.city,
    state       = EXCLUDED.state,
    postal_code = EXCLUDED.postal_code,
    country     = EXCLUDED.country,
    updated_at  = NOW()
RETURNING id, user_id, line1, line2, city, state, postal_code, country;

-- name: GetAddressByUserID :one
SELECT id, user_id, line1, line2, city, state, postal_code, country
FROM addresses WHERE user_id = $1;

-- name: SetActivationToken :exec
UPDATE users
SET activation_token        = $1,
    activation_expires      = $2,
    activation_requested_at = NOW(),
    last_email_sent_at          = NOW(),
    email_send_count            = resend_count + 1
WHERE id = $3;

-- name: GetUserByActivationToken :one
SELECT id, email, is_active, activation_token, activation_expires
FROM users WHERE activation_token = $1;

-- name: ActivateUser :exec
UPDATE users
SET is_active          = true,
    activation_token   = NULL,
    activation_expires = NULL
WHERE id = $1;

-- name: GetUserActivationMeta :one
SELECT id, email, is_active, email_send_count, last_email_sent_at
FROM users WHERE email = $1;