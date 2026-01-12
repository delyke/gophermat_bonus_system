package repository

import (
	"context"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

type BonusRepository interface {
	GetUserByLoginAndPassword(ctx context.Context, login, password string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}
