package repository

import (
    "context"
    "database/sql"
    "time"
    "user-api/db/sqlc"
)

type AuthRepository struct {
    q *sqlc.Queries
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
    return &AuthRepository{q: sqlc.New(db)}
}

func (r *AuthRepository) CreateUser(ctx context.Context, name, email, passwordHash, role string, dob time.Time) (sqlc.CreateUserRow, error) {
    return r.q.CreateUser(ctx, sqlc.CreateUserParams{
        Name:         name,
        Email:        email,
        PasswordHash: passwordHash,
        Role:         role,
        Dob:          dob,
    })
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
    return r.q.GetUserByEmail(ctx, email)
}

func (r *AuthRepository) SetActivationToken(ctx context.Context, token string, expires time.Time, userID int32) error {
    return r.q.SetActivationToken(ctx, sqlc.SetActivationTokenParams{
        ActivationToken:   sql.NullString{String: token, Valid: true},
        ActivationExpires: sql.NullTime{Time: expires, Valid: true},
        ID:                userID,
    })
}

func (r *AuthRepository) GetUserByActivationToken(ctx context.Context, token string) (sqlc.GetUserByActivationTokenRow, error) {
    return r.q.GetUserByActivationToken(ctx, sql.NullString{String: token, Valid: true})
}

func (r *AuthRepository) ActivateUser(ctx context.Context, userID int32) error {
    return r.q.ActivateUser(ctx, userID)
}

func (r *AuthRepository) GetUserActivationMeta(ctx context.Context, email string) (sqlc.GetUserActivationMetaRow, error) {
    return r.q.GetUserActivationMeta(ctx, email)
}