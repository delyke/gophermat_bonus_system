package v1

import (
	"errors"
	"strings"

	"github.com/delyke/gophermat_bonus_system/internal/model"
	bonusV1 "github.com/delyke/gophermat_bonus_system/pkg/openapi/bonus/v1"
)

func (s *ApiSuite) TestLoginUserBadCredentials() {
	req := &bonusV1.LoginRequest{Login: "max", Password: ""}
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("Login", s.ctx, "max", "").Return("", model.ErrBadCredentials)

	res, err := s.api.LoginUser(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.LoginUserBadRequest{}, res)
}

func (s *ApiSuite) TestLoginUserUnauthorized() {
	req := &bonusV1.LoginRequest{Login: "max", Password: "123"}
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("Login", s.ctx, "max", "123").Return("", model.ErrUserNotFound)

	res, err := s.api.LoginUser(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.LoginUserUnauthorized{}, res)
}

func (s *ApiSuite) TestLoginUserInternalServerError() {
	req := &bonusV1.LoginRequest{Login: "max", Password: "123"}
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("Login", s.ctx, "max", "123").Return("", errors.New("unexpected error"))

	res, err := s.api.LoginUser(s.ctx, req)

	s.Require().NoError(err)
	s.Require().IsType(&bonusV1.LoginUserInternalServerError{}, res)
}

func (s *ApiSuite) TestLoginUserOK() {
	req := &bonusV1.LoginRequest{Login: "max", Password: "123"}
	s.bonusService.On("Users").Return(s.userService)
	s.userService.On("Login", s.ctx, "max", "123").Return("token", nil)

	res, err := s.api.LoginUser(s.ctx, req)

	s.Require().NoError(err)
	okRes, ok := res.(*bonusV1.LoginUserOK)
	s.Require().True(ok)
	s.Require().True(okRes.SetCookie.Set)
	s.Require().True(strings.Contains(okRes.SetCookie.Value, "access_token=token"))
}
