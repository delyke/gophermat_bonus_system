package user

import (
	"errors"

	"github.com/google/uuid"

	"github.com/delyke/gophermat_bonus_system/internal/authctx"
	"github.com/delyke/gophermat_bonus_system/internal/model"
)

func (s *ServiceSuite) TestGetBalanceUnauthorized() {
	_, _, err := s.service.GetBalance(s.ctx)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrUnauthorized)
}

func (s *ServiceSuite) TestGetBalanceCurrentBalanceError() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})
	unexpectedErr := errors.New("balance error")
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.userRepository.On("GetBalanceByUUID", s.ctx, userUUID).Return(0.0, unexpectedErr)

	_, _, err := s.service.GetBalance(s.ctx)

	s.Require().Error(err)
	s.Require().ErrorIs(err, unexpectedErr)
}

func (s *ServiceSuite) TestGetBalanceWithdrawalsBalanceError() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})
	unexpectedErr := errors.New("withdrawals error")
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.bonusRepository.On("Withdrawals").Return(s.withdrawalRepository)
	s.userRepository.On("GetBalanceByUUID", s.ctx, userUUID).Return(100.0, nil)
	s.withdrawalRepository.On("SumAmountByUserUUID", s.ctx, userUUID).Return(0.0, unexpectedErr)

	_, _, err := s.service.GetBalance(s.ctx)

	s.Require().Error(err)
	s.Require().ErrorIs(err, unexpectedErr)
}

func (s *ServiceSuite) TestGetBalanceSuccess() {
	userUUID := uuid.New()
	s.ctx = authctx.WithPrincipal(s.ctx, authctx.Principal{
		UserUUID: userUUID,
		Login:    "max",
	})
	s.bonusRepository.On("Users").Return(s.userRepository)
	s.bonusRepository.On("Withdrawals").Return(s.withdrawalRepository)
	s.userRepository.On("GetBalanceByUUID", s.ctx, userUUID).Return(100.0, nil)
	s.withdrawalRepository.On("SumAmountByUserUUID", s.ctx, userUUID).Return(10.0, nil)

	current, withdrawn, err := s.service.GetBalance(s.ctx)

	s.Require().NoError(err)
	s.Require().Equal(100.0, current)
	s.Require().Equal(10.0, withdrawn)
}
