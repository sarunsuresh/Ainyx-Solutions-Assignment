package handler

import (
    "errors"
    "user-api/internal/logger"
    "user-api/internal/models"
    "user-api/internal/service"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
    "go.uber.org/zap"
)

type AuthHandler struct {
    svc      *service.AuthService
    validate *validator.Validate
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
    return &AuthHandler{svc: svc, validate: validator.New()}
}

// POST /auth/signup
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
    var req models.SignupRequest
    if err := c.BodyParser(&req); err != nil {
        return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
    }
    if err := h.validate.Struct(req); err != nil {
        return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "validation error")
    }

    if err := h.svc.Signup(c.Context(), req); err != nil {
        if errors.Is(err, service.ErrEmailTaken) {
            return models.RespondError(c, fiber.StatusConflict, models.CodeConflict, "email already in use")
        }
        logger.Log.Error("signup failed", zap.Error(err))
        return models.RespondError(c, fiber.StatusInternalServerError, models.CodeServerError, "Sign up failed, server error")
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "signup successful"})
}

// POST /auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req models.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
    }
    if err := h.validate.Struct(req); err != nil {
        return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
    }

    token, err := h.svc.Login(c.Context(), req)
    if err != nil {
        if errors.Is(err, service.ErrInvalidCredentials) {
            return models.RespondError(c, fiber.StatusUnauthorized, models.CodeUnauthorized, "invalid credentials")
        }
        logger.Log.Error("login failed", zap.Error(err))
        return models.RespondError(c, fiber.StatusInternalServerError, models.CodeServerError, "invalid body")
    }

    return c.JSON(models.LoginResponse{Token: token})
}