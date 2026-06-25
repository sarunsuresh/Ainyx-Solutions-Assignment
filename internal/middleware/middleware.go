package middleware

import (
	"time"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"user-api/internal/logger"
	"user-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)


func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := uuid.NewString()
		c.Locals("requestID", id)
		c.Set("X-Request-ID", id)
		return c.Next()
	}
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next() 
		logger.Log.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", time.Since(start)),
			zap.Any("requestID", c.Locals("requestID")),
		)
		return err
	}
}

func RequireAuth(secret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            return models.RespondError(c, fiber.StatusUnauthorized, models.CodeUnauthorized, "invalid token")
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fiber.ErrUnauthorized
            }
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            return models.RespondError(c, fiber.StatusUnauthorized, models.CodeUnauthorized, "invalid token")
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Locals("user", models.AuthUser{
            ID:    int32(claims["user_id"].(float64)), 
            Email: claims["email"].(string),
            Role:  claims["role"].(string),
        })

        return c.Next()
    }
}

func RequireRole(role string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, ok := c.Locals("user").(models.AuthUser)
        if !ok || user.Role != role {
            return models.RespondError(c, fiber.StatusForbidden, models.CodeForbidden, "forbidden")
        }
        return c.Next()
    }
}
