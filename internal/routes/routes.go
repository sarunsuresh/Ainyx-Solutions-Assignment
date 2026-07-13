package routes

import (
	"database/sql"
	"time"
	"user-api/internal/handler"
	"user-api/internal/middleware"
	"user-api/internal/redis"
	"user-api/internal/websocket"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, h *handler.UserHandler, ah *handler.AuthHandler,hub *websocket.Hub,rdb *redis.Client,db *sql.DB, jwtSecret string) {
	// global middleware (runs on every request)
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())
	app.Use(middleware.Recovery())
	app.Get("/ws", websocket.Handler(hub))
	app.Get("/healthz", handler.HealthCheck(db, rdb))
	auth := app.Group("/auth")
	auth.Post("/signup", ah.Signup)
	auth.Post("/login", middleware.RateLimit(rdb, 5, 60*time.Second), ah.Login)

	
	users := app.Group("/users")
	users.Get("/", h.ListUsers)

	
	protected := app.Group("/users", middleware.RequireAuth(jwtSecret))
	protected.Use(middleware.RateLimit(rdb, 30, 60*time.Second))
	protected.Get("/me", h.GetMe)
	protected.Patch("/me/passwordupdate", h.PasswordUpdate)
	protected.Get("/:id", h.GetUser)
	protected.Put("/:id", h.UpdateUser)
	protected.Put("/:id/profile", h.UpdateProfile)

	
	admin := app.Group("/admin", middleware.RequireAuth(jwtSecret), middleware.RequireRole("admin"))
	admin.Get("/users", h.ListUsers)
	admin.Delete("/users/:id", h.DeleteUser)
}
