package order

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/luhn"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *service) Create(ctx context.Context, orderID []byte) (string, error) {
	s.logger.Debug(ctx, "readed order id:", zap.String("orderId", string(orderID)))
	isValid, err := luhn.Validate(orderID)
	if err != nil {
		s.logger.Error(ctx, "ошибка при валидации по алгоритму Луна", zap.Error(err))
		return "", model.ErrBadCredentials
	}
	if !isValid {
		s.logger.Debug(ctx, "номер заказа не валидный", zap.String("orderId", string(orderID)))
		return "", model.ErrOrderIDLuhnInvalid
	}

	principal, ok := authctx.PrincipalFrom(ctx)
	if !ok {
		return "", model.ErrUnauthorized
	}

	order := &model.Order{
		OrderID:    string(orderID),
		Status:     model.OrderNew,
		UserUUID:   principal.UserUUID,
		UploadedAt: time.Now(),
	}

	cOrder, err := s.bonusRepository.Orders().Create(ctx, order)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderIDAlreadyExists):
			cOrder, err = s.bonusRepository.Orders().GetByNumber(ctx, string(orderID))
			if err != nil {
				return "", err
			}
			if cOrder.UserUUID != principal.UserUUID {
				return "", model.ErrOrderBelongsToAnotherUser
			}
			return "", model.ErrOrderAlreadyUploaded
		default:
			return "", err
		}
	}

	s.accrualWorker.Enqueue(model.OrderJob{
		OrderUUID:   cOrder.UUID,
		OrderNumber: cOrder.OrderID,
		UserID:      cOrder.UserUUID,
		OrderStatus: cOrder.Status,
		Attempts:    0,
	})

	return cOrder.UUID.String(), nil
}
