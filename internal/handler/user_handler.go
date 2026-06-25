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

// GET /users/me — requires auth middleware
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	// read the user that auth middleware injected
	authUser, ok := c.Locals("user").(models.AuthUser)
	if !ok {
		return models.RespondError(c, fiber.StatusUnauthorized, models.CodeUnauthorized, "unauthorized")
	}

	user, err := h.svc.GetUser(c.Context(), authUser.ID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(user)
}

// GET /users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid id")
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
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid id")
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
	}
	if err := h.validate.Struct(req); err != nil {
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
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
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid id")
	}

	if err := h.svc.DeleteUser(c.Context(), id); err != nil {
		logger.Log.Error("delete user failed", zap.Error(err))
		return models.RespondError(c, fiber.StatusInternalServerError, models.CodeServerError, "invalid body")
	}

	logger.Log.Info("user deleted", zap.Int32("id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) PasswordUpdate(c *fiber.Ctx) error {
	authUser, ok := c.Locals("user").(models.AuthUser)
	if !ok {
		return models.RespondError(c, fiber.StatusUnauthorized, models.CodeUnauthorized, "unauthorized")
	}

	var req models.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
	}
	if err := h.validate.Struct(req); err != nil {
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
	}

	user, err := h.svc.UpdatePassword(c.Context(), authUser.ID, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(user)

}

// GET /users?limit=&offset=
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit := queryInt(c, "limit", 10)
	offset := queryInt(c, "offset", 0)

	users, err := h.svc.ListUsers(c.Context(), limit, offset)
	if err != nil {
		logger.Log.Error("list users failed", zap.Error(err))
		return models.RespondError(c, fiber.StatusInternalServerError, models.CodeServerError, "loading users failed")
	}
	return c.JSON(users)
}

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
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeNotFound, "user not found ")
	}
	logger.Log.Error("service error", zap.Error(err))
	return models.RespondError(c, fiber.StatusBadRequest, models.CodeServerError, "internal server error")
}
