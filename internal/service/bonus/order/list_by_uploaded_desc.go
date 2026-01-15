package order

import (
	"context"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *service) ListByUploadedDesc(ctx context.Context) ([]*model.Order, error) {
	principal, ok := authctx.PrincipalFrom(ctx)
	if !ok {
		return nil, model.ErrUnauthorized
	}
	orders, err := s.bonusRepository.Orders().GetListByUploadedDesc(ctx, principal.UserUUID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
