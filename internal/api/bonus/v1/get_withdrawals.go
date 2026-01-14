package v1

import (
	"context"
	"errors"

	"github.com/samber/lo"

	"github.com/delyke/gophermat_bonus_system/internal/api/converter"
	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (a *api) GetWithdrawals(ctx context.Context) (bonusV1.GetWithdrawalsRes, error) {
	withdrawals, err := a.bonusService.Withdrawals().GetList(ctx)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUnauthorized):
			return &bonusV1.GetWithdrawalsUnauthorized{}, nil
		default:
			return &bonusV1.GetWithdrawalsInternalServerError{}, nil
		}
	}
	if len(withdrawals) == 0 {
		return &bonusV1.GetWithdrawalsNoContent{}, nil
	}
	apiWithdrawals := converter.ServiceWithdrawalListToApi(withdrawals)
	return lo.ToPtr(bonusV1.GetWithdrawalsResponse(apiWithdrawals)), nil
}
