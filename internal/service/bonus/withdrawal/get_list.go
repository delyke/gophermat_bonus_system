package withdrawal

import (
	"context"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *service) GetList(ctx context.Context) ([]*model.Withdrawal, error) {
	principal, ok := authctx.PrincipalFrom(ctx)
	if !ok {
		return nil, model.ErrUnauthorized
	}
	withdrawals, err := s.bonusRepository.Withdrawals().ListByUserByDateDesc(ctx, principal.UserUUID)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
