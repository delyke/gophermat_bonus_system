package workers

import (
	"context"
	"time"

	"github.com/delyke/gophermat_bonus_system/internal/model"
)

type AccrualWorker interface {
	Enqueue(job model.OrderJob)
	EnqueueLater(job model.OrderJob, d time.Duration)
	Stop()
	Bootstrap(ctx context.Context) error
}
