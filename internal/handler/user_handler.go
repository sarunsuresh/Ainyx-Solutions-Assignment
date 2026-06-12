package handler

import (
	"errors"
	"strconv"

	"user-api/internal/logger"
	"user-api/internal/models"
	"user-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc      *service.UserService
	validate *validator.Validate
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc, validate: validator.New()}
}

// POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.svc.CreateUser(c.Context(), req)
	if err != nil {
		logger.Log.Error("create user failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
	}

	logger.Log.Info("user created", zap.Int32("id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GET /users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	user, err := h.svc.GetUser(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(user)
}

// PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.svc.UpdateUser(c.Context(), id, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	logger.Log.Info("user updated", zap.Int32("id", id))
	return c.JSON(user)
}

// DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.svc.DeleteUser(c.Context(), id); err != nil {
		logger.Log.Error("delete user failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete user"})
	}

	logger.Log.Info("user deleted", zap.Int32("id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

// GET /users?limit=&offset=
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit := queryInt(c, "limit", 10)
	offset := queryInt(c, "offset", 0)

	users, err := h.svc.ListUsers(c.Context(), limit, offset)
	if err != nil {
		logger.Log.Error("list users failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not list users"})
	}
	return c.JSON(users)
}

// --- small helpers ---

func parseID(c *fiber.Ctx) (int32, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}

func queryInt(c *fiber.Ctx, key string, fallback int32) int32 {
	if v := c.Query(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return int32(n)
		}
	}
	return fallback
}

func handleServiceError(c *fiber.Ctx, err error) error {
	if errors.Is(err, service.ErrUserNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	logger.Log.Error("service error", zap.Error(err))
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
}
