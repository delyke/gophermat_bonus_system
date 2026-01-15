package v1

import (
	"errors"
	"strings"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (s *ApiSuite) TestRegisterUserBadCredentials() {
	req := &bonusV1.RegisterRequest{Login: "", Password: ""}
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("Register", s.ctx, "", "").Return("", model.ErrBadCredentials)

	res, err := s.api.RegisterUser(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.RegisterUserBadRequest{}, res)
}

func (s *ApiSuite) TestRegisterUserConflict() {
	req := &bonusV1.RegisterRequest{Login: "max", Password: "123"}
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("Register", s.ctx, "max", "123").Return("", model.ErrLoginTaken)

	res, err := s.api.RegisterUser(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.RegisterUserConflict{}, res)
}

func (s *ApiSuite) TestRegisterUserInternalServerError() {
	req := &bonusV1.RegisterRequest{Login: "max", Password: "123"}
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("Register", s.ctx, "max", "123").Return("", errors.New("unexpected error"))

	res, err := s.api.RegisterUser(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.RegisterUserInternalServerError{}, res)
}

func (s *ApiSuite) TestRegisterUserOK() {
	req := &bonusV1.RegisterRequest{Login: "max", Password: "123"}
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("Register", s.ctx, "max", "123").Return("token", nil)

	res, err := s.api.RegisterUser(s.ctx, req)

	s.Require().NoError(err)
	okRes, ok := res.(*bonusV1.RegisterUserOK)
	s.Require().True(ok)
	s.Require().True(okRes.SetCookie.Set)
	s.Require().True(strings.Contains(okRes.SetCookie.Value, "access_token=token"))
}
