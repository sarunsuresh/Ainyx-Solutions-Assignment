package routes

import (
	"user-api/internal/handler"
	"user-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, h *handler.UserHandler, ah *handler.AuthHandler, jwtSecret string) {
	// global middleware (runs on every request)
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())

	// public auth routes — no token needed
	auth := app.Group("/auth")
	auth.Post("/signup", ah.Signup)
	auth.Post("/login", ah.Login)

	
	users := app.Group("/users")
	users.Get("/", h.ListUsers)

	
	protected := app.Group("/users", middleware.RequireAuth(jwtSecret))
	protected.Get("/me", h.GetMe)
	protected.Patch("/me/passwordupdate", h.PasswordUpdate)
	protected.Get("/:id", h.GetUser)
	protected.Put("/:id", h.UpdateUser)

	
	admin := app.Group("/admin", middleware.RequireAuth(jwtSecret), middleware.RequireRole("admin"))
	admin.Get("/users", h.ListUsers)
	admin.Delete("/users/:id", h.DeleteUser)
}
