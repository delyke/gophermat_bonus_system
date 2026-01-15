package withdrawal

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/logger"
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
		logger.Error(ctx, "ошибка при валидации по алгоритму Луна", zap.Error(err))
		return model.ErrBadCredentials
	}
	if !isValid {
		logger.Debug(ctx, "номер заказа не валидный", zap.String("orderId", string(orderID)))
		return model.ErrOrderIDLuhnInvalid
	}

	cBalance, err := s.bonusRepository.Users().GetBalanceByUUID(ctx, principal.UserUUID)
	if err != nil {
		return err
	}
	if cBalance < sum {
		return model.ErrNotEnoughBalance
	}

	withdrawal := &model.Withdrawal{
		UserUUID:    principal.UserUUID,
		OrderID:     string(orderID),
		Amount:      sum,
		ProcessedAt: time.Now(),
	}

	err = s.bonusRepository.Withdrawals().Create(ctx, withdrawal)
	if err != nil {
		return err
	}
	newBalance := cBalance - sum
	err = s.bonusRepository.Users().SetBalanceByUUID(ctx, principal.UserUUID, newBalance)
	if err != nil {
		return err
	}
	return nil
}
