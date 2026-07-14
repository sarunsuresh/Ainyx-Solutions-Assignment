package main

import (
    "net"
    "log"

    "google.golang.org/grpc"
    "user-api/config"
    "user-api/internal/email"
    pb "user-api/internal/email/proto"
    "user-api/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
    _ = godotenv.Load()
    logger.Init()

    cfg := config.Load()

    lis, err := net.Listen("tcp", ":"+cfg.EmailServicePort)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterEmailServiceServer(grpcServer, email.NewServer(cfg.EmailFailureRate))

    logger.Log.Sugar().Infof("email service listening on :%s", cfg.EmailServicePort)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}