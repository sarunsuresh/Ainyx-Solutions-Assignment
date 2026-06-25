package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"user-api/config"
	"user-api/internal/models"
	"user-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrEmailTaken = errors.New("email already registered")

type AuthService struct {
	repo *repository.AuthRepository
	cfg  config.Config
}

func NewAuthService(repo *repository.AuthRepository, cfg config.Config) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

func (s *AuthService) Signup(ctx context.Context, req models.SignupRequest) error {
	_, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil {
		return ErrEmailTaken
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	dob, _ := time.Parse("2006-01-02", req.DOB)
	_, err = s.repo.CreateUser(ctx, req.Name, req.Email, string(hash), "user", dob)
	return err
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (string, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(s.cfg.JWTExpHours) * time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.cfg.JWTSecret))
}
