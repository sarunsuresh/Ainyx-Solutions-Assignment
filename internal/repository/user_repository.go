package repository

import (
	"context"
	"database/sql"
	"time"

	"user-api/db/sqlc"
)

// UserRepository handles all database access for users.
// It wraps the sqlc-generated Queries so the rest of the app
// never touches sqlc directly.
type UserRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{q: sqlc.New(db)}
}

func (r *UserRepository) Create(ctx context.Context, name string, dob time.Time) (sqlc.User, error) {
	return r.q.CreateUser(ctx, sqlc.CreateUserParams{Name: name, Dob: dob})
}

func (r *UserRepository) GetByID(ctx context.Context, id int32) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *UserRepository) Update(ctx context.Context, id int32, name string, dob time.Time) (sqlc.User, error) {
	return r.q.UpdateUser(ctx, sqlc.UpdateUserParams{ID: id, Name: name, Dob: dob})
}

func (r *UserRepository) Delete(ctx context.Context, id int32) error {
	return r.q.DeleteUser(ctx, id)
}

func (r *UserRepository) List(ctx context.Context, limit, offset int32) ([]sqlc.User, error) {
	return r.q.ListUsers(ctx, sqlc.ListUsersParams{Limit: limit, Offset: offset})
}
