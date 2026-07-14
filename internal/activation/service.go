package activation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"user-api/internal/clients"
	"user-api/internal/logger"
	"user-api/internal/queue"
	"user-api/internal/repository"

	"go.uber.org/zap"
)

var (
	ErrAlreadyActive  = errors.New("account already active")
	ErrCooldownActive = errors.New("resend not allowed yet")
	ErrInvalidToken   = errors.New("invalid or expired token")
	ErrUserNotFound   = errors.New("user not found")
)

const tokenExpiry = 24 * time.Hour

type ActivationService struct {
	repo        *repository.AuthRepository
	emailClient *clients.EmailClient
	queue       *queue.Queue
}

func NewActivationService(repo *repository.AuthRepository, emailClient *clients.EmailClient , queue *queue.Queue) *ActivationService {
	return &ActivationService{repo: repo, emailClient: emailClient, queue: queue }
}

func (s *ActivationService) SendActivation(ctx context.Context, userID int32, email string) error {
	token, err := generateToken()
	if err != nil {
		return err
	}
	expires := time.Now().Add(tokenExpiry)

	if err := s.repo.SetActivationToken(ctx, token, expires, userID); err != nil {
		logger.Log.Warn("email send failed, adding to retry queue",
			zap.String("email", email),
			zap.Error(err),
		)
		s.queue.Push(userID, email, token)
		return nil
	}

	
	return s.emailClient.SendActivationEmail(ctx, userID, email, token)
}

func (s *ActivationService) Activate(ctx context.Context, email, token string) error {
	user, err := s.repo.GetUserByActivationToken(ctx, token)
	if err != nil {
		return ErrInvalidToken
	}
	if user.IsActive {
		return ErrAlreadyActive
	}
	if time.Now().After(user.ActivationExpires.Time) {
		return ErrInvalidToken // expired
	}
	return s.repo.ActivateUser(ctx, user.ID)
}

func (s *ActivationService) Resend(ctx context.Context, email string) error {
	meta, err := s.repo.GetUserActivationMeta(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}
	if meta.IsActive {
		return ErrAlreadyActive
	}

	lastResent := meta.LastEmailSentAt.Time
	if !CanResend(int(meta.EmailSendCount), lastResent) {
		return ErrCooldownActive
	}

	return s.SendActivation(ctx, meta.ID, meta.Email)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
