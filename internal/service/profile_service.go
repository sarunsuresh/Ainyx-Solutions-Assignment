package service

import (
    "context"
    "database/sql"
    "time"

    "user-api/db/sqlc"
    "user-api/internal/models"
    "user-api/internal/repository"
)

type ProfileService struct {
    db      *sql.DB                      // needed to begin transactions
    userRepo *repository.UserRepository
    addrRepo *repository.AddressRepository
}

func NewProfileService(db *sql.DB, userRepo *repository.UserRepository, addrRepo *repository.AddressRepository) *ProfileService {
    return &ProfileService{db: db, userRepo: userRepo, addrRepo: addrRepo}
}

func (s *ProfileService) UpdateProfile (ctx context.Context , userID int32 , req models.UpdateProfileRequest) (models.ProfileResponse, error){
	dob, _ := time.Parse("2006-01-02", req.DOB)

	tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return models.ProfileResponse{}, err
    }

	defer tx.Rollback()

	user, err := s.userRepo.UpdateTx(ctx, tx, userID, req.Name,  dob)
    if err != nil {
        return models.ProfileResponse{}, err  

    }

	addr , err := s.addrRepo.UpsertTx(ctx, tx, sqlc.UpsertAddressParams{
        UserID:     userID,
        Line1:      req.Address.Line1,
        Line2:      sql.NullString{String: req.Address.Line2, Valid: req.Address.Line2 != ""},
        City:       req.Address.City,
        State:      req.Address.State,
        PostalCode: req.Address.PostalCode,
        Country:    req.Address.Country,
    })

	 if err != nil {
        return models.ProfileResponse{}, err  
    }

	if err := tx.Commit(); err != nil {
        return models.ProfileResponse{}, err
    }

	 return models.ProfileResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        DOB:   user.Dob.Format("2006-01-02"),
        Age:   CalculateAge(user.Dob),
        Address: models.AddressResponse{
            Line1:      addr.Line1,
            Line2:      addr.Line2.String,
            City:       addr.City,
            State:      addr.State,
            PostalCode: addr.PostalCode,
            Country:    addr.Country,
        },
    }, nil
}


