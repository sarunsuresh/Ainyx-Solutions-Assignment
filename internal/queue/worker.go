package queue

import (
    "context"
    "time"
    "user-api/internal/clients"
    "user-api/internal/logger"
    "go.uber.org/zap"
)

const (
    workerInterval  = 10 * time.Second
    initialBatch    = 5
    maxBatch        = 20
)

// Worker processes the email queue in the background.
type Worker struct {
    queue       *Queue
    emailClient *clients.EmailClient
    batchSize   int
}

func NewWorker(queue *Queue, emailClient *clients.EmailClient) *Worker {
    return &Worker{
        queue:       queue,
        emailClient: emailClient,
        batchSize:   initialBatch,
    }
}


func (w *Worker) Start(ctx context.Context) {
    go func() {
        ticker := time.NewTicker(workerInterval)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                logger.Log.Info("queue worker shutting down")
                return
            case <-ticker.C:
                w.processCycle(ctx)
            }
        }
    }()
    logger.Log.Info("queue worker started",
        zap.Duration("interval", workerInterval),
        zap.Int("batch_size", w.batchSize),
    )
}

func (w *Worker) processCycle(ctx context.Context) {
    if w.queue.Size() == 0 {
        return
    }

    // check if email service is healthy before draining
    ready, err := w.emailClient.CheckHealth(ctx)
    if err != nil || !ready {
        logger.Log.Warn("email service not ready, skipping queue drain",
            zap.Error(err),
            zap.Bool("ready", ready),
        )
        return
    }

    batch := w.queue.PopBatch(w.batchSize)
    if len(batch) == 0 {
        return
    }

    logger.Log.Info("processing queue batch",
        zap.Int("batch_size", len(batch)),
        zap.Int("queue_remaining", w.queue.Size()),
    )

    allSucceeded := true
    for _, task := range batch {
        err := w.emailClient.SendActivationEmail(
            ctx,
            task.UserID,
            task.Email,
            task.ActivationToken,
        )
        if err != nil {
            allSucceeded = false
            logger.Log.Error("queued task failed",
                zap.String("email", task.Email),
                zap.Int("retry_count", task.RetryCount),
                zap.Error(err),
            )
            w.queue.Requeue(task)
        } else {
            logger.Log.Info("queued task succeeded",
                zap.String("email", task.Email),
                zap.Int("retry_count", task.RetryCount),
            )
        }
    }

    // slow-release: adjust batch size based on results
    w.adjustBatchSize(allSucceeded)
}


func (w *Worker) adjustBatchSize(allSucceeded bool) {
    if allSucceeded {
        // gradually increase throughput
        if w.batchSize < maxBatch {
            w.batchSize = min(w.batchSize*2, maxBatch)
            logger.Log.Info("increasing batch size",
                zap.Int("new_batch_size", w.batchSize),
            )
        }
    } else {
        // back to safe minimum on any failure
        w.batchSize = initialBatch
        logger.Log.Warn("reducing batch size after failure",
            zap.Int("new_batch_size", w.batchSize),
        )
    }
}

