package v1

import (
	"context"
	"errors"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (a *api) RegisterUser(ctx context.Context, req *bonusV1.RegisterRequest) (bonusV1.RegisterUserRes, error) {
	login := req.Login
	password := req.Password

	token, err := a.bonusService.Users().Register(ctx, login, password)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrBadCredentials):
			return &bonusV1.RegisterUserBadRequest{}, nil
		case errors.Is(err, model.ErrLoginTaken):
			return &bonusV1.RegisterUserConflict{}, nil
		default:
			return &bonusV1.RegisterUserInternalServerError{}, nil
		}
	}
	c := buildAccessCookie(token)
	return &bonusV1.RegisterUserOK{
		SetCookie: bonusV1.OptString{
			Value: c.String(),
			Set:   true,
		},
	}, nil
}
