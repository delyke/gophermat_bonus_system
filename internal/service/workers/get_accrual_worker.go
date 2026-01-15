package workers

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/api/accrual/converter"
	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/repository"
	accrualV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/accrual/v1"
)

type AccrualProcessor struct {
	jobs      chan model.OrderJob
	semaphore chan struct{}

	rateMu    sync.Mutex
	rateUntil time.Time

	repo          repository.BonusRepository
	accrualClient *accrualV1.Client

	stopped atomic.Bool
	wg      sync.WaitGroup
	cancel  context.CancelFunc

	ctx context.Context //nolint:containedctx
}

func NewAccrualProcessor(
	ctx context.Context,
	repo repository.BonusRepository,
	accrualClient *accrualV1.Client,
	workers int,
) *AccrualProcessor {
	ctx, cancel := context.WithCancel(ctx)

	p := &AccrualProcessor{
		jobs:          make(chan model.OrderJob, 1000),
		semaphore:     make(chan struct{}, workers),
		repo:          repo,
		accrualClient: accrualClient,
		cancel:        cancel,
		ctx:           ctx,
	}

	go p.run(ctx)

	return p
}

// Bootstrap - метод для поднятия, берет из БД заказы, зависшие в статусе ожидания
func (p *AccrualProcessor) Bootstrap(ctx context.Context) error {
	const limit = 1000

	jobs, err := p.repo.Orders().ListPending(ctx, limit)
	if err != nil {
		return err
	}

	for _, j := range jobs {
		job := model.OrderJob{
			UserID:      j.UserUUID,
			OrderUUID:   j.UUID,
			OrderNumber: j.OrderID,
			OrderStatus: j.Status,
			Attempts:    0,
		}

		select {
		case p.jobs <- job:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	logger.Info(ctx, fmt.Sprintf("Accrual Воркер поднялся, взято в работу %d заказов", len(jobs)), zap.Int("count", len(jobs)))
	return nil
}

// Enqueue - добавить заказ в очередь на обработку
func (p *AccrualProcessor) Enqueue(job model.OrderJob) {
	if p.stopped.Load() {
		return
	}
	p.jobs <- job
}

// EnqueueLater - добавить заказ в очередь на обработку спустя некоторое время
func (p *AccrualProcessor) EnqueueLater(job model.OrderJob, d time.Duration) {
	go func() {
		timer := time.NewTimer(d)
		defer timer.Stop()

		select {
		case <-timer.C:
			p.Enqueue(job)
		case <-p.ctx.Done():
			return
		}
	}()
}

// Stop - закрытие канала и завершение работы для GracefulShutdown
func (p *AccrualProcessor) Stop() {
	p.stopped.Store(true)
	p.cancel()
	p.wg.Wait()
}

// run - запускает очередь со стоп-краном
func (p *AccrualProcessor) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-p.jobs:
			if !ok {
				return
			}

			select {
			case p.semaphore <- struct{}{}: // если контекст не отменен, то пытаемся получить слот
			case <-ctx.Done(): // если контекст отменен - то выходим
				return
			}

			p.wg.Add(1)
			go func(j model.OrderJob) {
				defer p.wg.Done()
				defer func() { <-p.semaphore }()
				p.process(ctx, j)
			}(job)
		}
	}
}

// waitIfRateLimited - функция, которая замораживает воркеров, если во внешней системе rate limit
func (p *AccrualProcessor) waitIfRateLimited(ctx context.Context) {
	p.rateMu.Lock()
	until := p.rateUntil
	p.rateMu.Unlock()

	now := time.Now()
	if !now.Before(until) {
		return
	}

	sleep := time.Until(until)
	timer := time.NewTimer(sleep)
	defer timer.Stop()

	select {
	case <-timer.C:
		return
	case <-ctx.Done():
		return
	}
}

func (p *AccrualProcessor) process(ctx context.Context, j model.OrderJob) {
	p.waitIfRateLimited(ctx)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if j.OrderStatus != model.OrderProcessing {
		err := p.repo.Orders().SetStatusByUUID(ctx, j.OrderUUID, model.OrderProcessing)
		if err != nil {
			logger.Error(ctx, "Ошибка при установке статуса для заказа", zap.Error(err))
			return
		}
		j.OrderStatus = model.OrderProcessing
	}

	result, err := p.accrualClient.GetOrderInfo(ctx, accrualV1.GetOrderInfoParams{
		Number: j.OrderNumber,
	})
	if err != nil {
		logger.Error(ctx, "[ACCRUAL] Ошибка получения информации о заказе при получении - попробуем чуть позже", zap.Error(err))
		p.backoffWithJitter(&j)
		return
	}

	switch v := result.(type) {
	case *accrualV1.GetOrderInfoNoContent:
		p.orderNoContentProcessing(ctx, &j)
		return
	case *accrualV1.GetOrderInfoResponse:
		switch v.Status {
		case accrualV1.OrderStatusINVALID:
			p.orderInvalidResponseProcessing(ctx, &j, v)
			return
		case accrualV1.OrderStatusPROCESSED:
			p.orderProcessedResponseProcessing(ctx, &j, v)
			return
		case accrualV1.OrderStatusREGISTERED:
			p.orderRegisteredResponseProcessing(ctx, &j)
			return
		case accrualV1.OrderStatusPROCESSING:
			p.orderProcessingResponseProcessing(ctx, &j)
			return
		default:
			p.orderUndefinedResponseProcessing(ctx, &j, v)
			return
		}
	case *accrualV1.GetOrderInfoTooManyRequestsHeaders:
		p.tooManyRequestsHandler(ctx, &j, v)
		return
	case *accrualV1.GetOrderInfoInternalServerError:
		p.internalServerHandler(ctx, &j)
		return
	}
}

func (p *AccrualProcessor) internalServerHandler(ctx context.Context, j *model.OrderJob) {
	logger.Debug(ctx, "На сервере accrual произошла ошибка, отправляем в работу еще раз...")
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) tooManyRequestsHandler(ctx context.Context, j *model.OrderJob, v *accrualV1.GetOrderInfoTooManyRequestsHeaders) {
	retryAfter := time.Duration(v.RetryAfter.Value) * time.Second

	p.rateMu.Lock()
	until := time.Now().Add(retryAfter)
	if until.After(p.rateUntil) {
		p.rateUntil = until
	}
	p.rateMu.Unlock()

	logger.Warn(ctx, "Пришел Rate Limit, ставим на стоп воркеров",
		zap.Duration("retry_after", retryAfter))
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) backoffWithJitter(j *model.OrderJob) {
	j.Attempts++
	delay := backoff(j.Attempts)
	delay += jitter(500 * time.Millisecond)
	p.EnqueueLater(*j, delay)
}

func (p *AccrualProcessor) orderUndefinedResponseProcessing(ctx context.Context, j *model.OrderJob, v *accrualV1.GetOrderInfoResponse) {
	logger.Debug(ctx, "Пришел неизвестный статус от системы начисления...", zap.Any("status", v.Status))
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) orderProcessingResponseProcessing(ctx context.Context, j *model.OrderJob) {
	logger.Debug(ctx, "Расчет начисления в процессе, отправляем в работу еще раз немного позже...")
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) orderRegisteredResponseProcessing(ctx context.Context, j *model.OrderJob) {
	logger.Debug(ctx, "Заказ зарегистрирован, отправляем в работу еще раз...")
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) orderNoContentProcessing(ctx context.Context, j *model.OrderJob) {
	logger.Debug(ctx, "accrual: 204 no content, retry later")
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) orderInvalidResponseProcessing(ctx context.Context, j *model.OrderJob, v *accrualV1.GetOrderInfoResponse) {
	err := p.repo.Orders().SetStatusByUUID(ctx, j.OrderUUID, model.OrderInvalid)
	if err != nil {
		logger.Error(ctx, "Ошибка при установке статуса заказа", zap.Error(err))
		return
	}
	logger.Debug(ctx, fmt.Sprintf("Заказ %s обработан со статусом %s", j.OrderNumber, v.Status))
}

func (p *AccrualProcessor) orderProcessedResponseProcessing(ctx context.Context, j *model.OrderJob, v *accrualV1.GetOrderInfoResponse) {
	err := p.repo.Orders().SetStatusByUUID(ctx, j.OrderUUID, model.OrderProcessed)
	if err != nil {
		logger.Error(ctx, "Ошибка при установке статуса начисленного вознаграждения", zap.Error(err))
		return
	}
	accrual := converter.OptFloat64ToFloat64(v.Accrual)
	if accrual != nil {
		err = p.repo.Orders().SetAccrualByUUID(ctx, j.OrderUUID, *accrual)
		if err != nil {
			logger.Error(ctx, "Ошибка при установке вознаграждения", zap.Error(err))
			return
		}
	}
}

func jitter(max time.Duration) time.Duration {
	return time.Duration(time.Now().UnixNano() % int64(max))
}

func backoff(attempt int) time.Duration {
	d := time.Second * time.Duration(1<<attempt)
	if d > time.Minute {
		d = time.Minute
	}
	return d
}
