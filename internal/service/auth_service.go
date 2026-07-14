package service

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"user-api/config"
	"user-api/internal/activation"
	"user-api/internal/logger"
	"user-api/internal/models"
	"user-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrEmailTaken = errors.New("email already registered")
var ErrAccountNotActivated = errors.New("account not activated ")

type AuthService struct {
	repo          *repository.AuthRepository
	cfg           config.Config
	activationSvc *activation.ActivationService
}

func NewAuthService(repo *repository.AuthRepository, cfg config.Config, activationSvc *activation.ActivationService) *AuthService {
	return &AuthService{repo: repo, cfg: cfg, activationSvc: activationSvc}
}

func (s *AuthService) Signup(ctx context.Context, req models.SignupRequest) (models.UserResponse, error) {
	_, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil {
		return models.UserResponse{}, ErrEmailTaken
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return models.UserResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserResponse{}, err
	}

	dob, _ := time.Parse("2006-01-02", req.DOB)
	u, err := s.repo.CreateUser(ctx, req.Name, req.Email, string(hash), "user", dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	if err := s.activationSvc.SendActivation(ctx, u.ID, u.Email); err != nil {
		logger.Log.Error("activation email failed after signup",
			zap.Error(err),
			zap.Int32("user_id", u.ID),
		)
		// don't return err — user account exists, email can be resent
	}
	return models.UserResponse{ID: u.ID,
		Name:  u.Name,
		Email: u.Email,
		DOB:   u.Dob.Format(dateLayout)}, err
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

	if !user.IsActive {
		return "", ErrAccountNotActivated
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(s.cfg.JWTExpHours) * time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.cfg.JWTSecret))
}
