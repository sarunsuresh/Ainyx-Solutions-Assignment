package repository

import (
    "context"
    "database/sql"
    "time"

    "user-api/db/sqlc"
)

type UserRepository struct {
    q *sqlc.Queries
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{q: sqlc.New(db)}
}

// Create removed — user creation is now handled by AuthRepository

func (r *UserRepository) GetByID(ctx context.Context, id int32) (sqlc.GetUserByIdRow, error) {
    return r.q.GetUserById(ctx, id)
}

func (r *UserRepository) Update(ctx context.Context, id int32, name string, email string, dob time.Time) (sqlc.UpdateUserRow, error) {
    return r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
        ID:    id,
        Name:  name,
        Email: email,
        Dob:   dob,
    })
}

func (r *UserRepository) Delete(ctx context.Context, id int32) error {
    return r.q.DeleteUser(ctx, id)
}

func (r *UserRepository) List(ctx context.Context, limit, offset int32) ([]sqlc.ListUsersRow, error) {
    return r.q.ListUsers(ctx, sqlc.ListUsersParams{Limit: limit, Offset: offset})
}

func (r *UserRepository) UpdatePass(ctx context.Context,PasswordHash string, id int32) ( sqlc.UpdatePasswordRow,error){
	return  r.q.UpdatePassword(ctx,sqlc.UpdatePasswordParams{PasswordHash: PasswordHash,ID: id})
	
}