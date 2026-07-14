package handler

import (
	"errors"
	"user-api/internal/activation"
	"user-api/internal/logger"
	"user-api/internal/models"
	"user-api/internal/service"
	"user-api/internal/websocket"
	"user-api/internal/clients"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	svc      *service.AuthService
	validate *validator.Validate
	wsHub    *websocket.Hub
	activationSvc *activation.ActivationService
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

// POST /auth/activate
func (h *AuthHandler) Activate(c *fiber.Ctx) error {
    var req struct {
        Email string `json:"email" validate:"required,email"`
        Token string `json:"token" validate:"required"`
    }
    if err := c.BodyParser(&req); err != nil {
        return models.RespondError(c, 400, models.CodeBadRequest, "invalid body")
    }
    if err := h.validate.Struct(req); err != nil {
        return err
    }

    if err := h.activationSvc.Activate(c.Context(), req.Email, req.Token); err != nil {
        switch {
        case errors.Is(err, activation.ErrAlreadyActive):
            return models.RespondError(c, 409, "ALREADY_ACTIVE", "account already active")
        case errors.Is(err, activation.ErrInvalidToken):
            return models.RespondError(c, 400, "INVALID_TOKEN", "invalid or expired token")
        default:
            return models.RespondError(c, 500, models.CodeServerError, "activation failed")
        }
    }

    return c.JSON(fiber.Map{"message": "account activated"})
}

// POST /auth/resend-activation
func (h *AuthHandler) ResendActivation(c *fiber.Ctx) error {
    var req struct {
        Email string `json:"email" validate:"required,email"`
    }
    if err := c.BodyParser(&req); err != nil {
        return models.RespondError(c, 400, models.CodeBadRequest, "invalid body")
    }
    if err := h.validate.Struct(req); err != nil {
        return err
    }

    if err := h.activationSvc.Resend(c.Context(), req.Email); err != nil {
        switch {
        case errors.Is(err, activation.ErrAlreadyActive):
            return models.RespondError(c, 409, "ALREADY_ACTIVE", "account already active")
        case errors.Is(err, activation.ErrCooldownActive):
            return models.RespondError(c, 429, "COOLDOWN_ACTIVE", "please wait before resending")
        case errors.Is(err, activation.ErrUserNotFound):
            return models.RespondError(c, 404, models.CodeNotFound, "user not found")
        case errors.Is(err, clients.ErrCircuitOpen):
            return models.RespondError(c, 503, "EMAIL_UNAVAILABLE", "email service unavailable, try again later")
        default:
            return models.RespondError(c, 500, models.CodeServerError, "resend failed")
        }
    }

    return c.JSON(fiber.Map{"message": "activation email sent"})
}