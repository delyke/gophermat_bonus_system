package user

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *ServiceSuite) TestRegisterBadCredentials() {
	_, err := s.service.Register(s.ctx, "", "")

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrBadCredentials)
}

func (s *ServiceSuite) TestRegisterCreateUserError() {
	unexpectedErr := errors.New("create user error")
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.userRepository.On("Create", s.ctx, "max", mock.Anything).Return(nil, unexpectedErr)

	_, err := s.service.Register(s.ctx, "max", "Strong123")

	s.Require().Error(err)
	s.Require().ErrorIs(err, unexpectedErr)
}

func (s *ServiceSuite) TestRegisterTokenIssueError() {
	userUUID := uuid.New()
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.userRepository.On("Create", s.ctx, "max", mock.Anything).Return(&model.User{
		UUID:  userUUID,
		Login: "max",
	}, nil)
	issueErr := errors.New("issue token error")
	s.tokenIssuer.On("IssueAccessToken", userUUID.String(), "max", mock.Anything).Return("", issueErr)

	_, err := s.service.Register(s.ctx, "max", "Strong123")

	s.Require().Error(err)
	s.Require().ErrorIs(err, issueErr)
}

func (s *ServiceSuite) TestRegisterSuccess() {
	userUUID := uuid.New()
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.userRepository.On("Create", s.ctx, "max", mock.Anything).Return(&model.User{
		UUID:  userUUID,
		Login: "max",
	}, nil)
	s.tokenIssuer.On("IssueAccessToken", userUUID.String(), "max", mock.Anything).Return("token", nil)

	token, err := s.service.Register(s.ctx, "max", "Strong123")

	s.Require().NoError(err)
	s.Require().Equal("token", token)
}

func (s *ServiceSuite) TestRegisterPasswordNotComplex() {
	_, err := s.service.Register(s.ctx, "max", "weak")

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrBadCredentials)
}
