package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-api/config"
	"user-api/internal/activation"
	"user-api/internal/clients"
	"user-api/internal/handler"
	"user-api/internal/logger"
	"user-api/internal/middleware"
	"user-api/internal/queue"
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
	cb := clients.NewCircuitBreaker(cfg.CBMaxFailures, cfg.CBResetTimeout)

	emailClient, err := clients.NewEmailClient(cfg.EmailServiceAddr, cb)
	emailQueue := queue.NewQueue()
	worker     := queue.NewWorker(emailQueue, emailClient)
	authRepo := repository.NewAuthRepository(db)

	activationSvc := activation.NewActivationService(authRepo, emailClient, emailQueue)

	if err != nil {
		log.Fatalf("failed to connect to email service: %v", err)
	}

	workerCtx, workerCancel := context.WithCancel(context.Background())
	worker.Start(workerCtx)

	hub := websocket.NewHub()
	fmt.Printf("main hub: %p\n", hub)
	repo := repository.NewUserRepository(db)
	addrepo := repository.NewAddressRepository(db)
	svc := service.NewUserService(repo)
	psvc := service.NewProfileService(db, repo, addrepo)
	h := handler.NewUserHandler(svc, psvc, hub)

	authSvc := service.NewAuthService(authRepo, cfg, activationSvc)

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
	workerCancel() 
	rdb.Close()
	db.Close()
	logger.Log.Info("shutdown complete")
}
