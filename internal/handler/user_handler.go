package handler

import (
	"errors"
	"strconv"
	"database/sql"
	"user-api/internal/logger"
	"user-api/internal/models"
	"user-api/internal/service"
	"user-api/internal/websocket"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"user-api/internal/redis"

)

type UserHandler struct {
    svc        *service.UserService
    profileSvc *service.ProfileService  // add this
    validate   *validator.Validate
	wsHub		*websocket.Hub
}

func NewUserHandler(svc *service.UserService, profileSvc *service.ProfileService , hub *websocket.Hub) *UserHandler {
    return &UserHandler{
        svc:        svc,
        profileSvc: profileSvc,
        validate:   validator.New(),
		wsHub: hub,
    }
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
		return err
	}

	user, err := h.svc.UpdateUser(c.Context(), id, req)
	if err != nil {
		return err
	}

	logger.Log.Info("user updated", zap.Int32("id", id))
	h.wsHub.Broadcast("user_updated", id)
	return c.JSON(user)
}

// DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
	}

	if err := h.svc.DeleteUser(c.Context(), id); err != nil {
		logger.Log.Error("delete user failed", zap.Error(err))
		return err
	}

	logger.Log.Info("user deleted", zap.Int32("id", id))
	h.wsHub.Broadcast("user_deleted", id)
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
		return err
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
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) UpdateProfile (c *fiber.Ctx) error{
	var req models.UpdateProfileRequest

	if err:=c.BodyParser(&req); err!=nil{
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid body")
	}

	if err:= h.validate.Struct(req);err!=nil{
		return err
	}
	id,err:=parseID(c);
	if err!=nil{
		return models.RespondError(c, fiber.StatusBadRequest, models.CodeBadRequest, "invalid id")
	}


	user,err:=h.profileSvc.UpdateProfile(c.Context(),id,req)
	if err!=nil{
		
		return handleServiceError(c, err)
	}

	h.wsHub.Broadcast("user_profile_updated", id)
	return c.JSON(user)

	

	}


func HealthCheck(db *sql.DB, rdb *redis.Client) fiber.Handler {
    return func(c *fiber.Ctx) error {
        if err := db.PingContext(c.Context()); err != nil {
            return c.Status(503).JSON(fiber.Map{
                "status": "degraded",
                "postgres": "down",
            })
        }
        if err := rdb.Ping(c.Context()); err != nil {
            return c.Status(503).JSON(fiber.Map{
                "status": "degraded",
                "redis": "down",
            })
        }
        return c.JSON(fiber.Map{"status": "ok"})
    }
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
		return err
	}
	logger.Log.Error("service error", zap.Error(err))
	return err
}
