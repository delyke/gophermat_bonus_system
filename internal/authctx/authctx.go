package authctx

import (
	"context"

	"github.com/google/uuid"
)

// Уникальный тип-ключ
type key struct{}

// Principal - identity пользователя после успешной аутентификации
type Principal struct {
	UserUUID uuid.UUID
	Login    string
}

// WithPrincipal - создает новый контекст с данными пользователя
func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, key{}, p)
}

// PrincipalFrom - безопасно достает пользователя из контекста
func PrincipalFrom(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(key{}).(Principal)
	return p, ok
}
