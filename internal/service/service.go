package service

import (
	"context"
	"time"
)

type BonusService interface {
	Users() UserService
}

type UserService interface {
	Login(ctx context.Context, login string, password string) (string, error)
	Register(ctx context.Context, login string, password string) (string, error)
}

type TokenIssuer interface {
	IssueAccessToken(userUUID, login string, ttl time.Duration) (string, error)
}
