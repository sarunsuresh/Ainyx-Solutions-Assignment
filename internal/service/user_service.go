package service

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
	"user-api/internal/models"
	"user-api/internal/repository"
)

// ErrUserNotFound is returned when a user doesn't exist.
// The handler maps this to a 404.
var ErrUserNotFound = errors.New("user not found")

const dateLayout = "2006-01-02"

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}
	return age
}


func (s *UserService) GetUser(ctx context.Context, id int32) (models.UserWithAgeResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserWithAgeResponse{}, ErrUserNotFound
		}
		return models.UserWithAgeResponse{}, err
	}
	return models.UserWithAgeResponse{ID: u.ID,
		Name: u.Name,
		DOB:  u.Dob.Format(dateLayout),
		Email: u.Email,
		Age:  CalculateAge(u.Dob),
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
    dob, _ := time.Parse(dateLayout, req.DOB)
    u, err := s.repo.Update(ctx, id, req.Name, req.Email, dob) 
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return models.UserResponse{}, ErrUserNotFound
        }
        return models.UserResponse{}, err
    }
    return models.UserResponse{
        ID:    u.ID,
        Name:  u.Name,
        Email: u.Email, 
        DOB:   u.Dob.Format(dateLayout),
    }, nil
}

func (s * UserService) UpdatePassword(ctx context.Context,id int32 , req models.UpdatePasswordRequest) (models.UserResponse,error){
	_ , err:= s.repo.GetByID(ctx,id)
	if err!=nil{
		return models.UserResponse{}, ErrUserNotFound
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return models.UserResponse{},err
    }

	data, err :=s.repo.UpdatePass(ctx,string(hash),id)

	if err != nil {
        return models.UserResponse{},err
    }

	return models.UserResponse{
		 ID:    data.ID,
        Name:  data.Name,
        Email: data.Email, 
        DOB:   data.Dob.Format(dateLayout),

	},nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int32) ([]models.UserWithAgeResponse, error) {
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]models.UserWithAgeResponse, 0, len(users))
	for _, u := range users {
		out = append(out, models.UserWithAgeResponse{
			ID:   u.ID,
			Name: u.Name,
			DOB:  u.Dob.Format(dateLayout),
			Age:  CalculateAge(u.Dob),
		})
	}
	return out, nil
}
