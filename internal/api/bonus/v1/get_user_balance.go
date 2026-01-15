package v1

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (a *api) GetUserBalance(ctx context.Context) (bonusV1.GetUserBalanceRes, error) {
	current, withdrawn, err := a.bonusService.Users().GetBalance(ctx)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUnauthorized):
			return &bonusV1.GetUserBalanceUnauthorized{}, nil
		default:
			return &bonusV1.GetUserBalanceInternalServerError{}, nil
		}
	}
	return &bonusV1.BalanceResponse{
		Current:   current,
		Withdrawn: withdrawn,
	}, nil
}
