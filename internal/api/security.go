package api

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

type SecurityHandler struct {
	secret []byte
}

// New - инициализирует SecurityHandler с заданным секретом для подписания JWT
func New(secret string) *SecurityHandler {
	return &SecurityHandler{secret: []byte(secret)}
}

// HandleJwtCookieAuth - проверяет наличие токена при запросе и валидирует его.
// Если на каком-то этапе произошла ошибка - отсекает запрос
// В случае если токен валиден - возвращает контекст с пользователем
func (s *SecurityHandler) HandleJwtCookieAuth(
	ctx context.Context,
	operationName bonusV1.OperationName,
	t bonusV1.JwtCookieAuth,
) (context.Context, error) {
	tokenStr := t.APIKey
	if tokenStr == "" {
		return ctx, ogenerrors.ErrSecurityRequirementIsNotSatisfied
	}

	principal, err := s.validateHS256(tokenStr)
	if err != nil {
		return ctx, ogenerrors.ErrSecurityRequirementIsNotSatisfied
	}

	return authctx.WithPrincipal(ctx, principal), nil
}

func (s *SecurityHandler) validateHS256(tokenStr string) (authctx.Principal, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !tok.Valid {
		return authctx.Principal{}, errors.New("invalid token")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return authctx.Principal{}, errors.New("invalid claims")
	}

	sub, _ := claims["sub"].(string)
	login, _ := claims["login"].(string)

	if sub == "" {
		return authctx.Principal{}, errors.New("missing sub")
	}

	userUUID, err := uuid.Parse(sub)
	if err != nil {
		return authctx.Principal{}, errors.New("invalid sub")
	}

	return authctx.Principal{
		UserUUID: userUUID,
		Login:    login,
	}, nil
}
