package v1

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (a *api) BalanceWithdrawal(ctx context.Context, req *bonusV1.BalanceWithdrawalRequest) (bonusV1.BalanceWithdrawalRes, error) {
	orderID := req.Order
	sum := req.Sum

	err := a.bonusService.Withdrawals().Create(ctx, []byte(orderID), sum)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUnauthorized):
			return &bonusV1.BalanceWithdrawalUnauthorized{}, nil
		case errors.Is(err, model.ErrNotEnoughBalance):
			return &bonusV1.BalanceWithdrawalPaymentRequired{}, nil
		case errors.Is(err, model.ErrOrderIDLuhnInvalid):
			return &bonusV1.BalanceWithdrawalUnprocessableEntity{}, nil
		case errors.Is(err, model.ErrBadCredentials):
			return &bonusV1.BalanceWithdrawalUnprocessableEntity{}, nil
		default:
			return &bonusV1.BalanceWithdrawalInternalServerError{}, nil
		}
	}
	return &bonusV1.BalanceWithdrawalOK{}, err
}
