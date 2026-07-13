package repository

import (
    "context"
    "database/sql"
    "user-api/db/sqlc"
)

type AddressRepository struct {
    db *sql.DB   
}               

func NewAddressRepository(db *sql.DB) *AddressRepository {
    return &AddressRepository{db: db}
}


func (r *AddressRepository) UpsertTx(ctx context.Context, tx *sql.Tx, arg sqlc.UpsertAddressParams) (sqlc.UpsertAddressRow, error) {
    q := sqlc.New(tx)   
    return q.UpsertAddress(ctx, arg)
}

func (r *AddressRepository) GetByUserID(ctx context.Context, userID int32) (sqlc.GetAddressByUserIDRow, error) {
    q := sqlc.New(r.db)
    return q.GetAddressByUserID(ctx, userID)
}