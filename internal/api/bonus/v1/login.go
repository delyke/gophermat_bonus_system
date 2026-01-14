package v1

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (a *api) LoginUser(ctx context.Context, req *bonusV1.LoginRequest) (bonusV1.LoginUserRes, error) {
	login := req.Login
	password := req.Password

	token, err := a.bonusService.Users().Login(ctx, login, password)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrBadCredentials):
			return &bonusV1.LoginUserBadRequest{}, nil
		case errors.Is(err, model.ErrUserNotFound):
			return &bonusV1.LoginUserUnauthorized{}, nil
		default:
			return &bonusV1.LoginUserInternalServerError{}, nil
		}
	}

	c := buildAccessCookie(token)
	return &bonusV1.LoginUserOK{
		SetCookie: bonusV1.OptString{
			Value: c.String(),
			Set:   true,
		},
	}, nil
}
