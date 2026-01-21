package v1

import (
	"bytes"
	"context"
	"io"

	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (a *api) OrderNumberLoad(ctx context.Context, req bonusV1.OrderNumberLoadReq) (bonusV1.OrderNumberLoadRes, error) {
	raw, err := io.ReadAll(req.Data)
	if err != nil {
		a.logger.Error(ctx, "Ошибка при чтении body", zap.Error(err))
		return &bonusV1.OrderNumberLoadBadRequest{}, nil
	}

	orderID := bytes.TrimSpace(raw)
	if len(orderID) == 0 {
		return &bonusV1.OrderNumberLoadBadRequest{}, nil
	}

	_, err = a.bonusService.Orders().Create(ctx, orderID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrBadCredentials):
			return &bonusV1.OrderNumberLoadBadRequest{}, nil
		case errors.Is(err, model.ErrOrderIDLuhnInvalid):
			return &bonusV1.OrderNumberLoadUnprocessableEntity{}, nil
		case errors.Is(err, model.ErrUnauthorized):
			return &bonusV1.OrderNumberLoadUnauthorized{}, nil
		case errors.Is(err, model.ErrOrderBelongsToAnotherUser):
			return &bonusV1.OrderNumberLoadConflict{}, nil
		case errors.Is(err, model.ErrOrderAlreadyUploaded):
			return &bonusV1.OrderNumberLoadOK{}, nil
		default:
			return &bonusV1.OrderNumberLoadInternalServerError{}, nil
		}
	}

	return &bonusV1.OrderNumberLoadAccepted{}, nil
}
