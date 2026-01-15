package user

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *ServiceSuite) TestLoginBadCredentials() {
	_, err := s.service.Login(s.ctx, "", "")

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrBadCredentials)
}

func (s *ServiceSuite) TestLoginUserNotFound() {
	unexpectedErr := errors.New("user not found")
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.userRepository.On("GetByLogin", s.ctx, "max").Return(nil, unexpectedErr)

	_, err := s.service.Login(s.ctx, "max", "123")

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrUserNotFound)
}

func (s *ServiceSuite) TestLoginPasswordMismatch() {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	s.Require().NoError(err)
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.userRepository.On("GetByLogin", s.ctx, "max").Return(&model.User{
		UUID:     uuid.New(),
		Login:    "max",
		Password: string(passwordHash),
	}, nil)

	_, err = s.service.Login(s.ctx, "max", "wrong")

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrUserNotFound)
}

func (s *ServiceSuite) TestLoginTokenIssueError() {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	s.Require().NoError(err)
	userUUID := uuid.New()
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.userRepository.On("GetByLogin", s.ctx, "max").Return(&model.User{
		UUID:     userUUID,
		Login:    "max",
		Password: string(passwordHash),
	}, nil)
	issueErr := errors.New("issue token error")
	s.tokenIssuer.On("IssueAccessToken", userUUID.String(), "max", mock.Anything).Return("", issueErr)

	_, err = s.service.Login(s.ctx, "max", "correct")

	s.Require().Error(err)
	s.Require().ErrorIs(err, issueErr)
}

func (s *ServiceSuite) TestLoginSuccess() {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	s.Require().NoError(err)
	userUUID := uuid.New()
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.userRepository.On("GetByLogin", s.ctx, "max").Return(&model.User{
		UUID:     userUUID,
		Login:    "max",
		Password: string(passwordHash),
	}, nil)
	s.tokenIssuer.On("IssueAccessToken", userUUID.String(), "max", mock.Anything).Return("token", nil)

	token, err := s.service.Login(s.ctx, "max", "correct")

	s.Require().NoError(err)
	s.Require().Equal("token", token)
}
