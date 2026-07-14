package clients

import (
	"errors"
	"testing"
	"time"
	"user-api/internal/logger"
)

func TestCircuitBreakerClosedToOpen(t *testing.T) {
	logger.Init()
	cb := NewCircuitBreaker(3, 30*time.Second)
	fail := func() error { return errors.New("fail") }

	cb.Call(fail)
	cb.Call(fail)
	cb.Call(fail)

	err := cb.Call(fail)
	if !errors.Is(err, ErrCircuitOpen) {
		t.Errorf("expected circuit open error, got %v", err)
	}
}

func TestCircuitBreakerHalfOpen(t *testing.T) {
	logger.Init()
	cb := NewCircuitBreaker(1, 10*time.Millisecond) // tiny timeout for test
	fail := func() error { return errors.New("fail") }
	succeed := func() error { return nil }

	cb.Call(fail) // opens circuit

	time.Sleep(20 * time.Millisecond) // wait for reset timeout

	// should allow one test call (half-open)
	err := cb.Call(succeed)
	if err != nil {
		t.Errorf("expected success in half-open, got %v", err)
	}

	// circuit should be closed again
	if cb.state != cbClosed {
		t.Error("circuit should be closed after successful half-open call")
	}
}
