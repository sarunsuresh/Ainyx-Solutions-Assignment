package main

import (
	"database/sql"
	"log"

	"user-api/config"
	"user-api/internal/handler"
	"user-api/internal/logger"
	"user-api/internal/repository"
	"user-api/internal/routes"
	"user-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	if err := logger.Init(); err != nil {
		log.Fatalf("logger init: %v", err)
	}
	defer logger.Log.Sync()

	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)
	authRepo := repository.NewAuthRepository(db)
	authSvc  := service.NewAuthService(authRepo, cfg)
	ah       := handler.NewAuthHandler(authSvc)


	app := fiber.New()
	routes.Register(app, h, ah, cfg.JWTSecret)

	logger.Log.Info(
	"server starting",
	zap.String("port", cfg.ServerPort),
)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}

