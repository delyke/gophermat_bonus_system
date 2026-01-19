package workers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/delyke/gophermat_bonus_system/internal/logger"
	"github.com/delyke/gophermat_bonus_system/internal/model"
	"github.com/delyke/gophermat_bonus_system/internal/repository/mocks"
	accrualV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/accrual/v1"
)

func newAccrualClient(t *testing.T, handler http.HandlerFunc) *accrualV1.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := accrualV1.NewClient(server.URL)
	require.NoError(t, err)

	return client
}

func TestNewAccrualProcessorStop(t *testing.T) {
	logger.SetNopLogger()
	repo := mocks.NewBonusRepository(t)
	client := newAccrualClient(t, http.NotFound)

	worker := NewAccrualProcessor(context.Background(), repo, client, 1)
	worker.Stop()
}

func TestBootstrapEnqueuesJobs(t *testing.T) {
	logger.SetNopLogger()
	repo := mocks.NewBonusRepository(t)
	orderRepo := mocks.NewOrderRepository(t)
	repo.On("Orders").Return(orderRepo)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	orderRepo.On("ListPending", mock.Anything, uint64(1000)).Return([]*model.Order{
		{
			UUID:     orderUUID,
			UserUUID: userUUID,
			OrderID:  "20000006",
			Status:   model.OrderNew,
		},
	}, nil)

	worker := &AccrualProcessor{
		jobs:          make(chan model.OrderJob, 1),
		semaphore:     make(chan struct{}, 1),
		repo:          repo,
		accrualClient: newAccrualClient(t, http.NotFound),
		ctx:           context.Background(),
	}

	err := worker.Bootstrap(context.Background())
	require.NoError(t, err)

	select {
	case job := <-worker.jobs:
		require.Equal(t, orderUUID, job.OrderUUID)
		require.Equal(t, userUUID, job.UserID)
		require.Equal(t, "20000006", job.OrderNumber)
		require.Equal(t, model.OrderNew, job.OrderStatus)
	default:
		t.Fatal("expected job to be enqueued")
	}
}

func TestEnqueueStoppedSkipsJob(t *testing.T) {
	worker := &AccrualProcessor{
		jobs: make(chan model.OrderJob, 1),
	}
	worker.stopped.Store(true)
	worker.Enqueue(model.OrderJob{OrderNumber: "20000006"})

	select {
	case <-worker.jobs:
		t.Fatal("expected no job to be enqueued")
	default:
	}
}

func TestEnqueueLaterEnqueuesJob(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker := &AccrualProcessor{
		jobs: make(chan model.OrderJob, 1),
		ctx:  ctx,
	}
	worker.EnqueueLater(model.OrderJob{OrderNumber: "20000006"}, 1*time.Millisecond)

	select {
	case job := <-worker.jobs:
		require.Equal(t, "20000006", job.OrderNumber)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected job to be enqueued later")
	}
}

func TestProcessInvalidResponse(t *testing.T) {
	logger.SetNopLogger()
	repo := mocks.NewBonusRepository(t)
	orderRepo := mocks.NewOrderRepository(t)
	repo.On("Orders").Return(orderRepo)

	orderUUID := uuid.New()
	orderRepo.On("SetStatusByUUID", mock.Anything, orderUUID, model.OrderProcessing).Return(nil)
	orderRepo.On("SetStatusByUUID", mock.Anything, orderUUID, model.OrderInvalid).Return(nil)

	client := newAccrualClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/orders/20000006", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprint(w, `{"order":"20000006","status":"INVALID"}`)
		if err != nil {
			return
		}
	})

	worker := &AccrualProcessor{
		jobs:          make(chan model.OrderJob, 1),
		semaphore:     make(chan struct{}, 1),
		repo:          repo,
		accrualClient: client,
		ctx:           context.Background(),
	}

	worker.process(context.Background(), model.OrderJob{
		OrderUUID:   orderUUID,
		OrderNumber: "20000006",
		OrderStatus: model.OrderNew,
	})

	orderRepo.AssertNotCalled(t, "SetAccrualByUUID", mock.Anything, mock.Anything, mock.Anything)
}

func TestProcessProcessedResponseWithAccrual(t *testing.T) {
	logger.SetNopLogger()
	repo := mocks.NewBonusRepository(t)
	orderRepo := mocks.NewOrderRepository(t)
	repo.On("Orders").Return(orderRepo)

	userRepo := mocks.NewUserRepository(t)
	repo.On("Users").Return(userRepo)

	orderUUID := uuid.New()
	orderRepo.On("SetStatusByUUID", mock.Anything, orderUUID, model.OrderProcessing).Return(nil)
	orderRepo.On("SetStatusByUUID", mock.Anything, orderUUID, model.OrderProcessed).Return(nil)
	orderRepo.On("SetAccrualByUUID", mock.Anything, orderUUID, 10.5).Return(nil)

	userRepo.On("GetBalanceByUUID", mock.Anything, mock.Anything).Return(10.5, nil)
	userRepo.On("SetBalanceByUUID", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	client := newAccrualClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/orders/20000006", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprint(w, `{"order":"20000006","status":"PROCESSED","accrual":10.5}`)
		if err != nil {
			return
		}
	})

	worker := &AccrualProcessor{
		jobs:          make(chan model.OrderJob, 1),
		semaphore:     make(chan struct{}, 1),
		repo:          repo,
		accrualClient: client,
		ctx:           context.Background(),
	}

	worker.process(context.Background(), model.OrderJob{
		OrderUUID:   orderUUID,
		OrderNumber: "20000006",
		OrderStatus: model.OrderNew,
	})
}

func TestBackoffWithJitterIncrementsAttempts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	worker := &AccrualProcessor{
		jobs: make(chan model.OrderJob, 1),
		ctx:  ctx,
	}

	job := model.OrderJob{Attempts: 0}
	worker.backoffWithJitter(&job)

	require.Equal(t, 1, job.Attempts)
}

func TestTooManyRequestsHandlerUpdatesRateLimit(t *testing.T) {
	err := logger.Init("debug", true)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	worker := &AccrualProcessor{
		jobs: make(chan model.OrderJob, 1),
		ctx:  ctx,
	}

	retry := 2
	resp := &accrualV1.GetOrderInfoTooManyRequestsHeaders{}
	resp.RetryAfter.SetTo(retry)

	now := time.Now()
	job := model.OrderJob{}
	worker.tooManyRequestsHandler(context.Background(), &job, resp)

	worker.rateMu.Lock()
	defer worker.rateMu.Unlock()
	require.True(t, worker.rateUntil.After(now.Add(time.Duration(retry-1)*time.Second)))
}

func TestInternalServerHandlerIncrementsAttempts(t *testing.T) {
	err := logger.Init("debug", true)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	worker := &AccrualProcessor{
		jobs: make(chan model.OrderJob, 1),
		ctx:  ctx,
	}

	job := model.OrderJob{}
	worker.internalServerHandler(context.Background(), &job)

	require.Equal(t, 1, job.Attempts)
}

func TestBackoffAndJitterHelpers(t *testing.T) {
	require.Equal(t, time.Second, backoff(0))
	require.Equal(t, 2*time.Second, backoff(1))
	require.Equal(t, time.Minute, backoff(10))

	j := jitter(500 * time.Millisecond)
	require.GreaterOrEqual(t, j, time.Duration(0))
	require.Less(t, j, 500*time.Millisecond)
}
