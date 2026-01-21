package withdrawal

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/luhn"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *service) Create(ctx context.Context, orderID []byte, sum float64) error {
	principal, ok := authctx.PrincipalFrom(ctx)
	if !ok {
		return model.ErrUnauthorized
	}
	isValid, err := luhn.Validate(orderID)
	if err != nil {
		s.logger.Error(ctx, "ошибка при валидации по алгоритму Луна", zap.Error(err))
		return model.ErrOrderIDLuhnInvalid
	}
	if !isValid {
		s.logger.Debug(ctx, "номер заказа не валидный", zap.String("orderId", string(orderID)))
		return model.ErrOrderIDLuhnInvalid
	}

	withdrawal := &model.Withdrawal{
		UserUUID:    principal.UserUUID,
		OrderID:     string(orderID),
		Amount:      sum,
		ProcessedAt: time.Now(),
	}

	err = s.bonusRepository.CreateWithdrawalWithBalance(ctx, withdrawal)
	if err != nil {
		return fmt.Errorf("failed to create user withdrawal: %w", err)
	}
	return nil
}
