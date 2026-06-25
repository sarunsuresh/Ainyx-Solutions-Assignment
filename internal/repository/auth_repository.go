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