package v1

import (
	"errors"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (s *ApiSuite) TestNewErrorUnauthorized() {
	errRes := s.api.NewError(s.ctx, ogenerrors.ErrSecurityRequirementIsNotSatisfied)

	s.Require().Equal(http.StatusUnauthorized, errRes.StatusCode)
	s.Require().Equal(bonusV1.NewOptInt(http.StatusUnauthorized), errRes.Response.Code)
}

func (s *ApiSuite) TestNewErrorInternalServerError() {
	errRes := s.api.NewError(s.ctx, errors.New("boom"))

	s.Require().Equal(http.StatusInternalServerError, errRes.StatusCode)
	s.Require().Equal(bonusV1.NewOptInt(http.StatusInternalServerError), errRes.Response.Code)
}
