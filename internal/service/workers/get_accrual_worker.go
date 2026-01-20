package workers

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
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

	rateUntil atomic.Int64

	pendingCursorMu   sync.Mutex
	pendingCursorTime time.Time
	pendingCursorUUID uuid.UUID

	repo          repository.BonusRepository
	accrualClient *accrualV1.Client
	logger        *logger.Logger

	stopped atomic.Bool
	wg      sync.WaitGroup
	cancel  context.CancelFunc

	ctx context.Context //nolint:containedctx
}

func (p *AccrualProcessor) retryAfterDuration(v *accrualV1.GetOrderInfoTooManyRequestsHeaders) time.Duration {
	if v == nil {
		return 0
	}
	return time.Duration(v.RetryAfter.Value) * time.Second
}

func NewAccrualProcessor(
	ctx context.Context,
	appLogger *logger.Logger,
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
		logger:        appLogger,
		cancel:        cancel,
		ctx:           ctx,
	}

	go p.run(ctx)

	return p
}

// Bootstrap - метод для поднятия, берет из БД заказы, зависшие в статусе ожидания
func (p *AccrualProcessor) Bootstrap(ctx context.Context) error {
	const limit = 1000

	orders, err := p.drainPending(ctx, limit, false)
	if err != nil {
		return err
	}

	go p.pollPending(ctx, limit, 2*time.Second)
	p.logger.Info(ctx, fmt.Sprintf("Accrual Воркер поднялся, взято в работу %d заказов", orders), zap.Int("count", orders))
	return nil
}

// Enqueue - добавить заказ в очередь на обработку
func (p *AccrualProcessor) Enqueue(job model.OrderJob) {
	if p.stopped.Load() {
		return
	}
	select {
	case p.jobs <- job:
	default:
		p.EnqueueLater(job, 100*time.Millisecond)
	}
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
	untilUnix := p.rateUntil.Load()
	if untilUnix == 0 {
		return
	}
	until := time.Unix(0, untilUnix)

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
			p.logger.Error(ctx, "Ошибка при установке статуса для заказа", zap.Error(err))
			return
		}
		j.OrderStatus = model.OrderProcessing
	}

	result, err := p.accrualClient.GetOrderInfo(ctx, accrualV1.GetOrderInfoParams{
		Number: j.OrderNumber,
	})
	if err != nil {
		p.logger.Error(ctx, "[ACCRUAL] Ошибка получения информации о заказе при получении - попробуем чуть позже", zap.Error(err))
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
	p.logger.Debug(ctx, "На сервере accrual произошла ошибка, отправляем в работу еще раз...")
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) tooManyRequestsHandler(ctx context.Context, j *model.OrderJob, v *accrualV1.GetOrderInfoTooManyRequestsHeaders) {
	retryAfter := p.retryAfterDuration(v)

	until := time.Now().Add(retryAfter)

	untilUnix := until.UnixNano()
	for {
		current := p.rateUntil.Load()
		if untilUnix <= current {
			break
		}
		if p.rateUntil.CompareAndSwap(current, untilUnix) {
			break
		}
	}

	p.logger.Warn(ctx, "Пришел Rate Limit, ставим на стоп воркеров",
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
	p.logger.Debug(ctx, "Пришел неизвестный статус от системы начисления...", zap.Any("status", v.Status))
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) orderProcessingResponseProcessing(ctx context.Context, j *model.OrderJob) {
	p.logger.Debug(ctx, "Расчет начисления в процессе, отправляем в работу еще раз немного позже...")
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) orderRegisteredResponseProcessing(ctx context.Context, j *model.OrderJob) {
	p.logger.Debug(ctx, "Заказ зарегистрирован, отправляем в работу еще раз...")
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) orderNoContentProcessing(ctx context.Context, j *model.OrderJob) {
	p.logger.Debug(ctx, "accrual: 204 no content, retry later")
	p.backoffWithJitter(j)
}

func (p *AccrualProcessor) orderInvalidResponseProcessing(ctx context.Context, j *model.OrderJob, v *accrualV1.GetOrderInfoResponse) {
	err := p.repo.Orders().SetStatusByUUID(ctx, j.OrderUUID, model.OrderInvalid)
	if err != nil {
		p.logger.Error(ctx, "Ошибка при установке статуса заказа", zap.Error(err))
		return
	}
	p.logger.Debug(ctx, fmt.Sprintf("Заказ %s обработан со статусом %s", j.OrderNumber, v.Status))
}

func (p *AccrualProcessor) orderProcessedResponseProcessing(ctx context.Context, j *model.OrderJob, v *accrualV1.GetOrderInfoResponse) {
	err := p.repo.Orders().SetStatusByUUID(ctx, j.OrderUUID, model.OrderProcessed)
	if err != nil {
		p.logger.Error(ctx, "Ошибка при установке статуса начисленного вознаграждения", zap.Error(err))
		return
	}
	accrual := converter.OptFloat64ToFloat64(v.Accrual)
	if accrual != nil {
		err = p.repo.Orders().SetAccrualByUUID(ctx, j.OrderUUID, *accrual)
		if err != nil {
			p.logger.Error(ctx, "Ошибка при установке вознаграждения", zap.Error(err))
			return
		}
	}
	balance, err := p.repo.Users().GetBalanceByUUID(ctx, j.UserID)
	if err != nil {
		p.logger.Error(ctx, "Ошибка при запросе баланса пользователя", zap.Error(err))
		return
	}
	balance += *accrual

	err = p.repo.Users().SetBalanceByUUID(ctx, j.UserID, balance)
	if err != nil {
		p.logger.Error(ctx, "Ошибка при установке баланса пользователю", zap.Error(err))
		return
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

func (p *AccrualProcessor) drainPending(ctx context.Context, limit uint64, useCursor bool) (int, error) {
	var (
		afterTime time.Time
		afterUUID uuid.UUID
		total     int
	)
	if useCursor {
		afterTime, afterUUID = p.getPendingCursor()
	}

	for {
		orders, err := p.repo.Orders().ListPendingAfter(ctx, limit, afterTime, afterUUID)
		if err != nil {
			return total, err
		}
		if len(orders) == 0 {
			break
		}

		for _, order := range orders {
			p.Enqueue(p.toOrderJob(order))
		}

		last := orders[len(orders)-1]
		afterTime = last.UploadedAt
		afterUUID = last.UUID
		p.setPendingCursor(afterTime, afterUUID)
		total += len(orders)
	}

	return total, nil
}

func (p *AccrualProcessor) pollPending(ctx context.Context, limit uint64, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if _, err := p.drainPending(ctx, limit, true); err != nil {
				p.logger.Error(ctx, "Ошибка при загрузке новых заказов", zap.Error(err))
			}
		}
	}
}

func (p *AccrualProcessor) setPendingCursor(afterTime time.Time, afterUUID uuid.UUID) {
	p.pendingCursorMu.Lock()
	p.pendingCursorTime = afterTime
	p.pendingCursorUUID = afterUUID
	p.pendingCursorMu.Unlock()
}

func (p *AccrualProcessor) getPendingCursor() (time.Time, uuid.UUID) {
	p.pendingCursorMu.Lock()
	defer p.pendingCursorMu.Unlock()
	return p.pendingCursorTime, p.pendingCursorUUID
}

func (p *AccrualProcessor) toOrderJob(order *model.Order) model.OrderJob {
	return model.OrderJob{
		UserID:      order.UserUUID,
		OrderUUID:   order.UUID,
		OrderNumber: order.OrderID,
		OrderStatus: order.Status,
		Attempts:    0,
	}
}
