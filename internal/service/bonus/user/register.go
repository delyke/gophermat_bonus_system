package user

import (
	"context"
	"github.com/delyke/gophermat_bonus_system/internal/config"
	"github.com/delyke/gophermat_bonus_system/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// Register - создает хэш пароля, записывает логин и пароль в БД и выдает токен
func (s *service) Register(ctx context.Context, login string, password string) (string, error) {
	if login == "" || password == "" {
		return "", model.ErrBadCredentials
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	created, err := s.bonusRepository.Users().Create(ctx, login, string(hashBytes))
	if err != nil {
		return "", err
	}

	tok, err := s.tokenIssuer.IssueAccessToken(
		created.UUID.String(),
		created.Login,
		config.Get().JWT.TTL(),
	)
	if err != nil {
		return "", err
	}

	return tok, nil
}
