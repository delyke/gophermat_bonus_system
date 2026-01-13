package v1

import (
	"context"
	"errors"
	"github.com/delyke/gophermat_bonus_system/internal/config"
	"github.com/delyke/gophermat_bonus_system/internal/service"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
	"github.com/ogen-go/ogen/ogenerrors"
	"net/http"
	"time"
)

type api struct {
	bonusV1.UnimplementedHandler
	bonusService service.BonusService
}

func NewApi(bs service.BonusService) *api {
	return &api{bonusService: bs}
}

func (a *api) NewError(_ context.Context, err error) *bonusV1.GenericErrorStatusCode {
	statusCode := http.StatusInternalServerError
	if err != nil {
		switch {
		case errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied):
			statusCode = http.StatusUnauthorized
		}
	}
	return &bonusV1.GenericErrorStatusCode{
		StatusCode: statusCode,
		Response: bonusV1.GenericError{
			Code:    bonusV1.NewOptInt(statusCode),
			Message: bonusV1.NewOptString(err.Error()),
		},
	}
}

func buildAccessCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  time.Now().Add(config.Get().JWT.TTL()),
	}
}
