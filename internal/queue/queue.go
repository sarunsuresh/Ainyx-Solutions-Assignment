package queue

import (
    "sync"
    "time"
    "user-api/internal/logger"
    "go.uber.org/zap"
)


type Queue struct {
    mu    sync.Mutex
    tasks []*Task
}

func NewQueue() *Queue {
    return &Queue{tasks: make([]*Task, 0)}
}

// Push adds a new task to the queue.
func (q *Queue) Push(userID int32, email, token string) {
    q.mu.Lock()
    defer q.mu.Unlock()

    task := &Task{
        UserID:          userID,
        Email:           email,
        ActivationToken: token,
        EnqueuedAt:      time.Now(),
        RetryCount:      0,
        NextRetryAt:     time.Now(), // due immediately on first attempt
    }
    q.tasks = append(q.tasks, task)

    logger.Log.Info("task queued",
        zap.String("email", email),
        zap.Int32("user_id", userID),
        zap.Int("queue_size", len(q.tasks)),
    )
}


func (q *Queue) PopBatch(n int) []*Task {
    q.mu.Lock()
    defer q.mu.Unlock()

    var batch []*Task
    var remaining []*Task

    for _, t := range q.tasks {
        if len(batch) < n && t.IsDue() {
            batch = append(batch, t)
        } else {
            remaining = append(remaining, t)
        }
    }

    q.tasks = remaining

    if len(batch) > 0 {
        logger.Log.Info("popped batch from queue",
            zap.Int("batch_size", len(batch)),
            zap.Int("queue_remaining", len(q.tasks)),
        )
    }

    return batch
}

// Requeue puts a failed task back with updated backoff.
func (q *Queue) Requeue(task *Task) {
    q.mu.Lock()
    defer q.mu.Unlock()

    task.RetryCount++
    task.NextRetryAt = time.Now().Add(NextBackoff(task.RetryCount))

    q.tasks = append(q.tasks, task)

    logger.Log.Warn("task requeued with backoff",
        zap.String("email", task.Email),
        zap.Int("retry_count", task.RetryCount),
        zap.Time("next_retry_at", task.NextRetryAt),
    )
}

// Size returns current number of queued tasks.
func (q *Queue) Size() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    return len(q.tasks)
}