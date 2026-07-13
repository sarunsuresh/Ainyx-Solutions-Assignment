package middleware

import(
	"database/sql"
	"errors"
	"fmt"
	"user-api/internal/logger"
	"user-api/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler(c *fiber.Ctx, err error) error {

	
	status := fiber.StatusInternalServerError
	message := "internal server error"
	var validationErr validator.ValidationErrors
	
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		status = fiberErr.Code
		message = fiberErr.Message
	}

	
	if errors.Is(err, sql.ErrNoRows) {
		status = fiber.StatusNotFound
		message = "resource not found"
	}

	if errors.As(err, &validationErr){
	status = fiber.StatusBadRequest

	field := validationErr[0]

	message = fmt.Sprintf(
		"%s failed validation (%s)",
		field.Field(),
		field.Tag(),
	)
}

if errors.Is(err, service.ErrEmailTaken) {
    status = fiber.StatusConflict
    message = "email already registered"
}

if errors.Is(err, service.ErrInvalidCredentials) {
    status = fiber.StatusBadRequest
    message = "email already registered"
}



	if status >= 500 {
		logger.Log.Error(
			"request failed",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
		)
	}

	return c.Status(status).JSON(fiber.Map{
		"error": message,
		"requestId": c.Locals("requestID"),
	})
}