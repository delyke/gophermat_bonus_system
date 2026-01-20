package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/delyke/gophermat_bonus_system/internal/model"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	Create(ctx context.Context, login, password string) (*model.User, error)
	GetBalanceByUUID(ctx context.Context, uuid uuid.UUID) (float64, error)
	SetBalanceByUUID(ctx context.Context, uuid uuid.UUID, amount float64) error
}

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
	GetByNumber(ctx context.Context, number string) (*model.Order, error)
	GetListByUploadedDesc(ctx context.Context, userUUID uuid.UUID) ([]*model.Order, error)
	SetStatusByUUID(ctx context.Context, uuid uuid.UUID, status model.OrderStatus) error
	SetAccrualByUUID(ctx context.Context, uuid uuid.UUID, accrual float64) error
	ListPending(ctx context.Context, limit uint64) ([]*model.Order, error)
	ListPendingAfter(ctx context.Context, limit uint64, afterUploadedAt time.Time, afterUUID uuid.UUID) ([]*model.Order, error)
}

type WithdrawalRepository interface {
	Create(ctx context.Context, withdrawal *model.Withdrawal) error
	ListByUserByDateDesc(ctx context.Context, userUUID uuid.UUID) ([]*model.Withdrawal, error)
	SumAmountByUserUUID(ctx context.Context, userUUID uuid.UUID) (float64, error)
}

type BonusRepository interface {
	Users() UserRepository
	Orders() OrderRepository
	Withdrawals() WithdrawalRepository
	CreateWithdrawalWithBalance(ctx context.Context, withdrawal *model.Withdrawal) error
}
