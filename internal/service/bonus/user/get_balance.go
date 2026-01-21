package user

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *service) GetBalance(ctx context.Context) (float64, float64, error) {
	principal, ok := authctx.PrincipalFrom(ctx)
	if !ok {
		return 0, 0, model.ErrUnauthorized
	}
	cBalance, err := s.bonusRepository.Users().GetBalanceByUUID(ctx, principal.UserUUID)
	if err != nil {
		s.logger.Error(ctx, "Ошибка при получении баланса пользователя", zap.Error(err))
		return 0, 0, err
	}
	wBalance, err := s.bonusRepository.Withdrawals().SumAmountByUserUUID(ctx, principal.UserUUID)
	if err != nil {
		s.logger.Error(ctx, "Ошибка при получении суммы выведенных средств", zap.Error(err))
		return 0, 0, fmt.Errorf("getting sum withdrawals by user error: %w", err)
	}
	return cBalance, wBalance, nil
}
