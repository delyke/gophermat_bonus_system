package service

import (
	"context"
	"time"

	"github.com/delyke/gophermat_bonus_system/internal/model"
)

type BonusService interface {
	Users() UserService
	Orders() OrderService
	Withdrawals() WithdrawalService
}

type UserService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, login, password string) (string, error)
	GetBalance(ctx context.Context) (float64, float64, error)
}

type OrderService interface {
	Create(ctx context.Context, orderID []byte) (string, error)
	GetByNumber(ctx context.Context, number string) (*model.Order, error)
	ListByUploadedDesc(ctx context.Context) ([]*model.Order, error)
}

type WithdrawalService interface {
	Create(ctx context.Context, orderID []byte, sum float64) error
	GetList(ctx context.Context) ([]*model.Withdrawal, error)
}

type TokenIssuer interface {
	IssueAccessToken(userUUID, login string, ttl time.Duration) (string, error)
}
