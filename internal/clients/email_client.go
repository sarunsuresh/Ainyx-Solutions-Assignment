package clients

import (
    "context"
    "errors"
    "sync"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "user-api/internal/email/proto"
    "user-api/internal/logger"
    "go.uber.org/zap"
)



type cbState int

const (
    cbClosed   cbState = iota  // normal
    cbOpen                     // failing fast
    cbHalfOpen                 // testing recovery
)

type CircuitBreaker struct {
    mu              sync.Mutex
    state           cbState
    failures        int
    maxFailures     int
    resetTimeout    time.Duration
    lastFailureTime time.Time
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures:  maxFailures,
        resetTimeout: resetTimeout,
    }
}

var ErrCircuitOpen = errors.New("circuit breaker open: email service unavailable")

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()

    switch cb.state {
    case cbOpen:
        if time.Since(cb.lastFailureTime) > cb.resetTimeout {
            cb.state = cbHalfOpen
            logger.Log.Info("circuit breaker half-open")
        } else {
            cb.mu.Unlock()
            logger.Log.Warn("circuit breaker open — fast fail")
            return ErrCircuitOpen
        }
    }

    cb.mu.Unlock()

    err := fn()

    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err != nil {
        cb.failures++
        cb.lastFailureTime = time.Now()
        if cb.failures >= cb.maxFailures {
            cb.state = cbOpen
            logger.Log.Error("circuit breaker opened",
                zap.Int("failures", cb.failures))
        }
        return err
    }

    if cb.state == cbHalfOpen {
        logger.Log.Info("circuit breaker closed — recovered")
    }
    cb.failures = 0
    cb.state = cbClosed
    return nil
}



type EmailClient struct {
    grpc *pb.EmailServiceClient  
    cb   *CircuitBreaker
}

func NewEmailClient(addr string, cb *CircuitBreaker) (*EmailClient, error) {
    conn, err := grpc.Dial(addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        return nil, err
    }
    client := pb.NewEmailServiceClient(conn)				
    return &EmailClient{grpc: &client, cb: cb}, nil
}

func (c *EmailClient) SendActivationEmail(ctx context.Context, userID int32, email, token string) error {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    return c.cb.Call(func() error {
        resp, err := (*c.grpc).SendActivationEmail(ctx, &pb.SendActivationEmailRequest{
            UserId:          userID,
            Email:           email,
            ActivationCode: token,
        })
        if err != nil {
            return err
        }
        if !resp.Success {
            return errors.New(resp.Messsage)
        }
        return nil
    })
}

func (c *EmailClient) CheckHealth(ctx context.Context) (bool, error) {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    resp, err := (*c.grpc).CheckHealth(ctx, &pb.HealthRequest{})
    if err != nil {
        return false, err
    }
    return resp.Ready, nil
}