package queue

import "time"

// Task represents a failed email send that needs to be retried.
type Task struct {
    UserID          int32
    Email           string
    ActivationToken string
    EnqueuedAt      time.Time
    RetryCount      int
    NextRetryAt     time.Time
}


var backoffDurations = []time.Duration{
    1 * time.Minute,
    2 * time.Minute,
    4 * time.Minute,
    8 * time.Minute,
    16 * time.Minute,
    30 * time.Minute,  // cap
}

// NextBackoff returns how long to wait before the next retry.
func NextBackoff(retryCount int) time.Duration {
    if retryCount >= len(backoffDurations) {
        return backoffDurations[len(backoffDurations)-1]
    }
    return backoffDurations[retryCount]
}

// IsDue returns true if the task is ready to be retried now.
func (t *Task) IsDue() bool {
    return time.Now().After(t.NextRetryAt)
}