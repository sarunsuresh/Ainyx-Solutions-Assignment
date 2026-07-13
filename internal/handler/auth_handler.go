package handler

import (
	"errors"
	"user-api/internal/logger"
	"user-api/internal/models"
	"user-api/internal/service"
	"user-api/internal/websocket"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	svc      *service.AuthService
	validate *validator.Validate
	wsHub    *websocket.Hub
}

func NewAuthHandler(svc *service.AuthService, hub *websocket.Hub) *AuthHandler {
	return &AuthHandler{svc: svc, validate: validator.New(), wsHub: hub}
}

// POST /auth/signup
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req models.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
	}
	if err := h.validate.Struct(req); err != nil {
		return err
	}


	u, err := h.svc.Signup(c.Context(), req)
	if err != nil {

		if errors.Is(err, service.ErrEmailTaken) {
			return models.RespondError(c, fiber.StatusConflict, models.CodeConflict, "email already in use")
		}

		logger.Log.Error("signup failed", zap.Error(err))
		return err
	}
	
	h.wsHub.Broadcast("user_created", u.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "signup successful",
	})
}

// POST /auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		 return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
	}
	if err := h.validate.Struct(req); err != nil {
		return err
	}

	token, err := h.svc.Login(c.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return models.RespondError(c, fiber.StatusUnauthorized, models.CodeUnauthorized, "invalid credentials")
		}
		logger.Log.Error("login failed", zap.Error(err))
		return err
	}

	return c.JSON(models.LoginResponse{Token: token})
}
