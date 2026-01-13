package repository

import (
	"context"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	Create(ctx context.Context, login, password string) (*model.User, error)
}

type BonusRepository interface {
	Users() UserRepository
}
