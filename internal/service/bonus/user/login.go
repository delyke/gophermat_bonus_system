package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/delyke/gophermat_bonus_system/internal/config"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

// Login - Проводит аутентификацию по паре логин / пароль
func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	if username == "" && password == "" {
		return "", model.ErrBadCredentials
	}

	user, err := s.bonusRepository.Users().GetByLogin(ctx, username)
	if err != nil {
		return "", model.ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return "", model.ErrUserNotFound
	}

	tok, err := s.tokenIssuer.IssueAccessToken(
		user.UUID.String(),
		user.Login,
		config.Get().JWT.TTL(),
	)
	if err != nil {
		return "", err
	}

	return tok, nil
}
