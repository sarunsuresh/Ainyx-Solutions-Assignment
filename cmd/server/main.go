package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"fmt"
	"os/signal"
	"syscall"
	"time"
	"user-api/config"
	"user-api/internal/handler"
	"user-api/internal/logger"
	"user-api/internal/middleware"
	"user-api/internal/redis"
	"user-api/internal/repository"
	"user-api/internal/routes"
	"user-api/internal/service"
	"user-api/internal/websocket"

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
	hub := websocket.NewHub()
	fmt.Printf("main hub: %p\n", hub)
	repo := repository.NewUserRepository(db)
	addrepo := repository.NewAddressRepository(db)
	svc := service.NewUserService(repo)
	psvc := service.NewProfileService(db, repo, addrepo)
	h := handler.NewUserHandler(svc, psvc, hub)
	authRepo := repository.NewAuthRepository(db)
	authSvc := service.NewAuthService(authRepo, cfg)

	ah := handler.NewAuthHandler(authSvc, hub)
	rdb, _ := redis.New(cfg.RedisURL)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	routes.Register(app, h, ah, hub, rdb, db, cfg.JWTSecret)

	logger.Log.Info(
		"server starting",
		zap.String("port", cfg.ServerPort),
	)

	go func() {
		if err := app.Listen(":" + cfg.ServerPort); err != nil {
			logger.Log.Error("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Log.Error("shutdown error", zap.Error(err))
	}

	rdb.Close()
	db.Close()
	logger.Log.Info("shutdown complete")
}
