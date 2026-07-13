package middleware

import (
	"time"
	"strings"
	"fmt"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"user-api/internal/logger"
	"user-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"user-api/internal/redis"
	"runtime/debug"

   
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





func RateLimit(rdb *redis.Client, maxRequests int, window time.Duration) fiber.Handler {
    return func(c *fiber.Ctx) error {
        ip := c.IP()
        key := fmt.Sprintf("rate:%s:%s", c.Path(), ip)

        count, err := rdb.Increment(context.Background(), key, window)
        if err != nil {
            logger.Log.Error("redis rate limit failed", zap.Error(err))
            return c.Next()
        }

        if count > int64(maxRequests) {
            c.Set("Retry-After", fmt.Sprintf("%d", int(window.Seconds())))
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "error":       "too many requests",
                "retry_after": int(window.Seconds()),
            })
        }

        c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
        c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", maxRequests-int(count)))

        return c.Next()
    }
}

func Recovery() fiber.Handler {
    return func(c *fiber.Ctx) error {
        defer func() {
            if r := recover(); r != nil {
                logger.Log.Error("panic recovered",
                    zap.Any("error", r),
                    zap.String("stack", string(debug.Stack())),
                    zap.String("path", c.Path()),
                )
                c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                    "error":     "internal server error",
                    "requestId": c.Locals("requestID"),
                })
            }
        }()
        return c.Next()
    }
}

