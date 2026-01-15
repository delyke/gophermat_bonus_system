package v1

import (
	"errors"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (s *APISuite) TestGetUserBalanceUnauthorized() {
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("GetBalance", s.ctx).Return(0.0, 0.0, model.ErrUnauthorized)

	res, err := s.api.GetUserBalance(s.ctx)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.GetUserBalanceUnauthorized{}, res)
}

func (s *APISuite) TestGetUserBalanceInternalServerError() {
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("GetBalance", s.ctx).Return(0.0, 0.0, errors.New("unexpected error"))

	res, err := s.api.GetUserBalance(s.ctx)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.GetUserBalanceInternalServerError{}, res)
}

func (s *APISuite) TestGetUserBalanceOK() {
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("GetBalance", s.ctx).Return(100.0, 10.0, nil)

	res, err := s.api.GetUserBalance(s.ctx)

	s.Require().NoError(err)
	okRes, ok := res.(*bonusV1.BalanceResponse)
	s.Require().True(ok)
	s.Require().Equal(100.0, okRes.Current)
	s.Require().Equal(10.0, okRes.Withdrawn)
}
