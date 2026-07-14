package email

import (
    "context"
    "math/rand"
    pb "user-api/internal/email/proto"
    "user-api/internal/logger"
    "go.uber.org/zap"
	"fmt"
)

type Server struct{
	pb.UnimplementedEmailServiceServer
	failureRate float64
}

func NewServer(failureRate float64) *Server{
	return &Server{failureRate: failureRate}
}

func (s *Server) SendActivationEmail(ctx context.Context,req *pb.SendActivationEmailRequest,)(*pb.SendActivationEmailResponse,error){
	if rand.Float64() < s.failureRate {
        logger.Log.Error("simulated email failure",
            zap.Int32("user_id", req.UserId),
            zap.String("email", req.Email),
        )
        return nil, fmt.Errorf("simulated email service failure")
    }

	 logger.Log.Info("activation email sent",
        zap.Int32("user_id", req.UserId),
        zap.String("email", req.Email),
        zap.String("token", req.ActivationCode),
        
    )

	return &pb.SendActivationEmailResponse{
        Success: true,
        Messsage: "email queued",
    }, nil
}

func (s *Server)CheckHealth(ctx context.Context , req *pb.HealthRequest)(*pb.HealthResponse,error){
	if s.failureRate >= 1.0 {
        return &pb.HealthResponse{
            Ready:   false,
            Message: "service not ready",
        }, nil
    }
    return &pb.HealthResponse{
        Ready:   true,
        Message: "ready",
    }, nil
}