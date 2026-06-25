package models

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
    Message   string `json:"message"`
    Code      string `json:"code"`
    RequestID string `json:"request_id"`
}

// ErrorCodes — named constants so you never typo a code string
const (
    CodeBadRequest   = "BAD_REQUEST"
    CodeUnauthorized = "UNAUTHORIZED"
    CodeForbidden    = "FORBIDDEN"
    CodeNotFound     = "NOT_FOUND"
    CodeConflict     = "CONFLICT"
    CodeServerError  = "INTERNAL_ERROR"
)


func RespondError(c *fiber.Ctx, status int, code, message string) error {
    requestID, _ := c.Locals("requestID").(string)
    return c.Status(status).JSON(ErrorResponse{
        Error: ErrorDetail{
            Message:   message,
            Code:      code,
            RequestID: requestID,
        },
    })
}